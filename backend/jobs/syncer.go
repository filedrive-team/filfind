package jobs

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filedrive-team/filfind/backend/filclient"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/pando"
	"github.com/filedrive-team/filfind/backend/repo"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	address.CurrentNetwork = address.Mainnet
}

type FilecoinSyncer struct {
	filClient *filclient.FilClient
	repo      *repo.Manager
}

func NewFilecoinSyncer(filClient *filclient.FilClient, repo *repo.Manager) *FilecoinSyncer {
	return &FilecoinSyncer{
		filClient: filClient,
		repo:      repo,
	}
}

func (s *FilecoinSyncer) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	halfHourTicker := time.NewTicker(30 * time.Minute)
	halfDayTicker := time.NewTicker(12 * time.Hour)
	//s.syncMinerOwner(ctx)
	//s.syncStorageDeal(ctx)
	s.syncMinerSuccessRate(ctx)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			halfHourTicker.Stop()
			return
		case <-ticker.C:
			s.syncMinerOwner(ctx)
			s.syncMinerFirstDealTimeFromLocal(ctx)
			s.syncClientDataCap(ctx)
			s.syncClientStatsFromLocal(ctx)
			s.syncClientSumDealsFromLocal(ctx)
		case <-halfHourTicker.C:
			s.syncReviewsStatsFromLocal(ctx)
		case <-halfDayTicker.C:
			s.syncProviderStatsMonthlyFromLocal(ctx)
			s.syncMinerSuccessRate(ctx)
		}
	}
}

func (s *FilecoinSyncer) Init(ctx context.Context) {
	s.syncMinerOwner(ctx)
	s.syncMinerSuccessRate(ctx)
	//s.syncStorageDeal(ctx)
	s.syncClientStatsFromLocal(ctx)
	s.syncProviderStatsMonthlyFromLocal(ctx)
}

func (s *FilecoinSyncer) syncMinerOwner(ctx context.Context) {
	logger.Info("sync miner owner...")
	defer logger.Info("sync miner owner finished")

	var mgr utils.AsyncManager
	limitCh := make(chan struct{}, 64)
	retryCh := make(chan address.Address)
	resCh := make(chan *models.Provider)
	defer func() {
		close(limitCh)
		close(retryCh)
		close(resCh)
		retryCh = nil
	}()
	limit := 100
	for offset := 0; ; offset += limit {
		list, err := s.repo.GetProviderAddress(limit, offset)
		if err != nil {
			logger.WithError(err).Error("call GetProviderAddress failed")
			break
		}
		if len(list) == 0 {
			break
		}
		updates := make([]*models.Provider, 0, limit)
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		task := func(addr address.Address) {
			limitCh <- struct{}{}
			defer func() {
				<-limitCh
			}()
			info, err := s.filClient.StateMinerInfo(ctx, addr, filclient.EmptyTSK)
			if err != nil {
				select {
				case retryCh <- addr:
				default:
				}
				return
			}
			resCh <- &models.Provider{
				Address: addr.String(),
				Owner:   info.Owner.String(),
			}
		}
		go func() {
			for _, miner := range list {
				addr, err := address.NewFromString(miner)
				if err != nil {
					logger.WithError(err).WithField("miner", miner).Error("call address.NewFromString failed")
					continue
				}
				mgr.AddTask(func() {
					task(addr)
				})
			}

			mgr.Wait()
			cancel()
		}()
		func() {
			for {
				select {
				case addr := <-retryCh:
					curAddr := addr
					mgr.AddTask(func() {
						task(curAddr)
					})
				case res := <-resCh:
					updates = append(updates, res)
				case <-ctx.Done():
					return
				}
			}
		}()

		err = s.repo.UpdateProviderOwner(updates)
		if err != nil {
			logger.WithError(err).Error("call UpdateProviderOwner failed")
		}
	}
}

func (s *FilecoinSyncer) syncStorageDeal(ctx context.Context) {
	logger.Info("sync storage deal from filecoin api...")
	defer logger.Info("sync storage deal from filecoin api finished")

	var mgr utils.AsyncManager
	limitCh := make(chan struct{}, 64)
	retryCh := make(chan uint64)
	resCh := make(chan *models.Deal)
	defer func() {
		close(limitCh)
		close(retryCh)
		close(resCh)
		retryCh = nil
	}()
	deal, err := s.repo.LastDeal()
	if err != nil {
		logger.WithError(err).Error("call LastDeal failed")
		return
	}
	offset := int64(1000)
	loops := 0
	for dealId := int64(deal.DealId); dealId >= 0; dealId -= offset {
		firstDealId := dealId - offset
		if firstDealId < 0 {
			firstDealId = 0
		}

		if loops%100 == 0 {
			count, err := s.repo.CountDeal(uint64(dealId))
			if err != nil {
				logger.WithError(err).Error("call CountDeal failed")
				break
			}
			if dealId+1 == count {
				break
			}
		}
		loops += 1

		count, err := s.repo.CountDealRange(uint64(firstDealId), uint64(dealId))
		if err != nil {
			logger.WithError(err).Error("call CountDealRange failed")
			break
		}
		if dealId-firstDealId == count {
			continue
		}

		updates := make([]*models.Deal, 0, offset)
		ctxTask, cancel := context.WithTimeout(ctx, 5*time.Minute)
		task := func(dealId uint64) {
			limitCh <- struct{}{}
			defer func() {
				<-limitCh
			}()
			ctx := context.TODO()
			//ts, err := s.filClient.ChainHead(ctx)
			//if err != nil {
			//	logger.WithError(err).Error("call ChainHead failed")
			//	return
			//}
			//ts, err = s.filClient.ChainGetTipSetByHeight(ctx, abi.ChainEpoch(1752558), ts.Key())
			//if err != nil {
			//	logger.WithError(err).Error("call ChainGetTipSetByHeight failed")
			//	return
			//}
			d, err := s.filClient.StateMarketStorageDeal(ctx, abi.DealID(dealId), filclient.EmptyTSK)
			if err != nil {
				if !strings.Contains(err.Error(), "not found") {
					select {
					case retryCh <- dealId:
					default:
					}
				}
				return
			}
			resCh <- &models.Deal{
				DealId:               abi.DealID(dealId),
				PieceCid:             d.Proposal.PieceCID.String(),
				PieceSize:            d.Proposal.PieceSize,
				VerifiedDeal:         d.Proposal.VerifiedDeal,
				Client:               d.Proposal.Client.String(),
				Provider:             d.Proposal.Provider.String(),
				Label:                d.Proposal.Label,
				StartEpoch:           d.Proposal.StartEpoch,
				EndEpoch:             d.Proposal.EndEpoch,
				StoragePricePerEpoch: decimal.NewFromBigInt(d.Proposal.StoragePricePerEpoch.Int, 0),
				ProviderCollateral:   decimal.NewFromBigInt(d.Proposal.ProviderCollateral.Int, 0),
				ClientCollateral:     decimal.NewFromBigInt(d.Proposal.ClientCollateral.Int, 0),
				SectorStartEpoch:     d.State.SectorStartEpoch,
				LastUpdatedEpoch:     d.State.LastUpdatedEpoch,
				SlashEpoch:           d.State.SlashEpoch,
			}
		}
		go func() {
			if dealIds, err := s.repo.GetDealIdByRange(uint64(firstDealId), uint64(dealId)); err == nil {
				index := 0
				for id := firstDealId; id < dealId; id += 1 {
					if index < len(dealIds) {
						if uint64(id) == dealIds[index] {
							index += 1
							continue
						}
					}
					curId := uint64(id)
					mgr.AddTask(func() {
						task(curId)
					})
				}
			}

			mgr.Wait()
			cancel()
		}()
		func() {
			for {
				select {
				case id := <-retryCh:
					curId := id
					mgr.AddTask(func() {
						task(curId)
					})
				case res := <-resCh:
					updates = append(updates, res)
				case <-ctxTask.Done():
					return
				}
			}
		}()

		err = s.repo.UpsertDeal(updates)
		if err != nil {
			logger.WithError(err).Error("call UpsertDeal failed")
		}
	}
}

func (s *FilecoinSyncer) syncMinerFirstDealTimeFromLocal(ctx context.Context) {
	logger.Info("sync miner first deal time...")
	defer logger.Info("sync miner first deal time finished")

	data, err := s.repo.StatsProviderFirstDealTime()
	if err != nil {
		logger.WithError(err).Error("call StatsProviderFirstDealTime failed")
		return
	}
	updates := make([]*models.Provider, 0, len(data))
	for _, deal := range data {
		updates = append(updates, &models.Provider{
			Address:       deal.Provider,
			FirstDealTime: utils.GetBlockTimeByEpoch(int64(deal.StartEpoch)).Unix(),
		})
	}
	err = s.repo.UpdateProviderFirstDealTime(updates)
	if err != nil {
		logger.WithError(err).Error("call UpdateProviderFirstDealTime failed")
	}
}

func (s *FilecoinSyncer) syncReviewsStatsFromLocal(ctx context.Context) {
	logger.Info("sync reviews stats...")
	defer logger.Info("sync reviews stats finished")

	data, err := s.repo.StatsReviews()
	if err != nil {
		logger.WithError(err).Error("call StatsReviews failed")
		return
	}
	updates := make([]*models.Provider, 0, len(data))
	for _, item := range data {
		updates = append(updates, &models.Provider{
			Address:     item.Provider,
			Reviews:     item.Reviews,
			ReviewScore: item.AvgScore,
		})
	}
	err = s.repo.UpdateProviderReviewsStats(updates)
	if err != nil {
		logger.WithError(err).Error("call UpdateProviderReviewsStats failed")
	}
}

func (s *FilecoinSyncer) syncClientDataCap(ctx context.Context) {
	logger.Info("sync client data cap from filecoin api...")
	defer logger.Info("sync client data cap from filecoin api finished")

	var mgr utils.AsyncManager
	limitCh := make(chan struct{}, 64)
	retryCh := make(chan *repo.ClientAddress)
	resCh := make(chan *models.ClientInfo)
	defer func() {
		close(limitCh)
		close(retryCh)
		close(resCh)
		retryCh = nil
	}()
	list, err := s.repo.ClientUserList()
	if err != nil {
		logger.WithError(err).Error("call ClientUserList failed")
		return
	}
	updates := make([]*models.ClientInfo, 0, len(list))
	ctxTask, cancel := context.WithTimeout(ctx, 5*time.Minute)
	task := func(cli *repo.ClientAddress) {
		limitCh <- struct{}{}
		defer func() {
			<-limitCh
		}()

		addr, err := address.NewFromString(cli.AddressId)
		if err != nil {
			logger.WithError(err).WithField("addrId", cli.AddressId).Error("call address.NewFromString failed")
			return
		}
		vcs, err := s.filClient.StateVerifiedClientStatus(ctxTask, addr, filclient.EmptyTSK)
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				select {
				case retryCh <- cli:
				default:
				}
			}
			return
		}
		resCh <- &models.ClientInfo{
			Uid:       cli.Uid,
			AddressId: addr.String(),
			DataCap:   decimal.NewFromBigInt(vcs.Int, 0),
		}
	}
	go func() {
		for _, cli := range list {
			curCli := cli
			mgr.AddTask(func() {
				task(curCli)
			})
		}

		mgr.Wait()
		cancel()
	}()
	func() {
		for {
			select {
			case cli := <-retryCh:
				curCli := cli
				mgr.AddTask(func() {
					task(curCli)
				})
			case res := <-resCh:
				updates = append(updates, res)
			case <-ctxTask.Done():
				return
			}
		}
	}()

	err = s.repo.UpsertClientInfoDataCap(updates)
	if err != nil {
		logger.WithError(err).Error("call UpsertClientInfoDataCap failed")
	}
}

func (s *FilecoinSyncer) syncClientStatsFromLocal(ctx context.Context) {
	logger.Info("sync client stats...")
	defer logger.Info("sync client stats finished")

	deal, err := s.repo.LastDeal()
	if err != nil {
		logger.WithError(err).Error("call LastDeal failed")
		return
	}
	sectionDealId := int64(0)
	lastClientStats, err := s.repo.LastClientStats()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.WithError(err).Error("call StatsClientVerifiedDeal failed")
			return
		}
	} else {
		// recalculate the previous stage
		sectionDealId = int64(lastClientStats.SectionDealId) - models.SectionDeals
		if sectionDealId < 0 {
			sectionDealId = 0
		}
	}
	for ; sectionDealId < int64(deal.DealId); sectionDealId += models.SectionDeals {
		list, err := s.repo.StatsClientDeal(sectionDealId, sectionDealId+models.SectionDeals)
		if err != nil {
			logger.WithError(err).Error("call StatsClientDeal failed")
			return
		}

		for _, item := range list {
			item.SectionDealId = abi.DealID(sectionDealId)
		}
		err = s.repo.UpsertClientStats(list)
		if err != nil {
			logger.WithError(err).Error("call UpsertClientStats failed")
		}
	}
}

func (s *FilecoinSyncer) syncClientSumDealsFromLocal(ctx context.Context) {
	logger.Info("sync client sum deals...")
	defer logger.Info("sync client sum deals finished")

	data, err := s.repo.GetClientStatsSum()
	if err != nil {
		logger.WithError(err).Error("call GetClientStatsSum failed")
		return
	}

	list := make([]*models.ClientInfo, 0, len(data))
	for _, item := range data {
		list = append(list, &models.ClientInfo{
			AddressId:       item.Client,
			StorageCapacity: item.StorageCapacity,
			StorageDeals:    item.StorageDeals,
			UsedDataCap:     item.UsedDataCap,
			VerifiedDeals:   item.VerifiedDeals,
		})
	}
	err = s.repo.UpsertClientInfoDeals(list)
	if err != nil {
		logger.WithError(err).Error("call UpsertClientInfoDeals failed")
	}

}

func (s *FilecoinSyncer) syncProviderStatsMonthlyFromLocal(ctx context.Context) {
	logger.Info("sync provider stats monthly...")
	defer logger.Info("sync provider stats monthly finished")

	var beginTime time.Time
	lastStats, err := s.repo.LastProviderStatsMonthly()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.WithError(err).Error("call LastProviderStatsMonthly failed")
			return
		} else {
			deal, err := s.repo.FirstValidDeal()
			if err != nil {
				logger.WithError(err).Error("call FirstValidDeal failed")
				return
			}
			beginTime = utils.GetBlockTimeByEpoch(int64(deal.StartEpoch))
		}
	} else {
		beginTime = lastStats.Month.Time.AddDate(0, -1, 0)
	}

	startMonth := utils.MonthBegin(beginTime)
	endMonth := utils.MonthBegin(time.Now().AddDate(0, 1, 0))
	var nextMonth time.Time
	for ; startMonth.Before(endMonth); startMonth = nextMonth {
		nextMonth = startMonth.AddDate(0, 1, 0)

		startEpoch := utils.GetEpochByTime(startMonth)
		endEpoch := utils.GetEpochByTime(nextMonth)
		list, err := s.repo.ProviderClientsByRange(startEpoch, endEpoch)
		if err != nil {
			logger.WithError(err).Error("call ProviderClientsByRange failed")
			return
		}
		psms := make([]*models.ProviderStatsMonthly, 0, len(list))
		for _, item := range list {
			psms = append(psms, &models.ProviderStatsMonthly{
				Month:    types.NewUnixTime(startMonth),
				Provider: item.Provider,
				Client:   item.Client,
			})
		}
		err = s.repo.UpsertProviderStatsMonthly(psms)
		if err != nil {
			logger.WithError(err).Error("call UpsertProviderStatsMonthly failed")
			return
		}
	}
}

func (s *FilecoinSyncer) syncMinerSuccessRate(ctx context.Context) {
	logger.Info("sync success rate from pando...")
	defer logger.Info("sync success rate from pando finished")

	keyFile := filepath.Join(os.TempDir(), "pandoPeerId.key")
	dtStats := pando.GetMinerDealTaskStats(keyFile)
	updates := make([]*models.Provider, 0, len(dtStats))
	for miner, stats := range dtStats {
		storageRate := decimal.Zero
		retrievalRate := decimal.Zero
		storageTotal := stats.Storage.Success + stats.Storage.Failed
		if storageTotal > 0 {
			storageRate = decimal.NewFromFloat(float64(stats.Storage.Success) / float64(storageTotal))
		}
		retrievalTotal := stats.Retrieval.Success + stats.Retrieval.Failed
		if retrievalTotal > 0 {
			retrievalRate = decimal.NewFromFloat(float64(stats.Retrieval.Success) / float64(retrievalTotal))
		}
		updates = append(updates, &models.Provider{
			Address:              miner,
			DealSuccessRate:      storageRate,
			RetrievalSuccessRate: retrievalRate,
		})
	}
	err := s.repo.UpdateProviderSuccessRate(updates)
	if err != nil {
		logger.WithError(err).Error("call UpdateProviderSuccessRate failed")
	}
}
