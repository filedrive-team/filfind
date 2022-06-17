package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filedrive-team/filfind/backend/filclient"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/repo"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type MinerList struct {
	Pagination *struct {
		Total  int
		Offset int
		Limit  int
	}
	Miners []*struct {
		Address         string
		Status          bool
		UptimeAverage   float64
		Price           *string
		VerifiedPrice   *string
		MinPieceSize    *decimal.Decimal
		MaxPieceSize    *decimal.Decimal
		RawPower        decimal.Decimal
		QualityAdjPower decimal.Decimal
		IsoCode         string
		Region          string
		Score           decimal.Decimal
		Scores          *struct {
			Total                  decimal.Decimal
			Uptime                 decimal.Decimal
			StorageDeals           decimal.Decimal
			CommittedSectorsProofs decimal.Decimal
		}
		FreeSpace    *decimal.Decimal
		StorageDeals *struct {
			Total           uint64
			NoPenalties     uint64
			SuccessRate     decimal.Decimal
			AveragePrice    decimal.Decimal
			DataStored      decimal.Decimal
			Slashed         uint64
			Terminated      uint64
			FaultTerminated uint64
		}
	}
}

type Crawler struct {
	filrepApi string
	repo      *repo.Manager
}

func NewCrawler(filrepApi string, repo *repo.Manager) *Crawler {
	return &Crawler{
		filrepApi: filrepApi,
		repo:      repo,
	}
}

func (c *Crawler) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	dealTicker := time.NewTicker(12 * time.Hour)
	//c.startFilrepTask(ctx)
	//c.startDealTask(ctx)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			dealTicker.Stop()
			return
		case <-ticker.C:
			c.startFilrepTask(ctx)
		case <-dealTicker.C:
			c.startDealTask(ctx)
		}
	}
}

func (c *Crawler) Init(ctx context.Context) {
	c.startFilrepTask(ctx)
	c.startDealTask(ctx)
}

func (c *Crawler) makeRequest(ctx context.Context, url string) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.19 Safari/537.36")
	return
}

func (c *Crawler) startFilrepTask(ctx context.Context) {
	logger.Info("sync filrep.io data...")
	defer logger.Info("sync filrep.io data finished")

	for page := 0; ; page += 1 {
		if err := c.syncMinerList(ctx, 100, page*100); err != nil {
			break
		}
	}
}

func (c *Crawler) minerList(ctx context.Context, url string) (*MinerList, error) {
	cli := &http.Client{}
	req, err := c.makeRequest(ctx, url)
	if err != nil {
		logger.WithError(err).WithField("url", url).Errorf("makeRequest failed")
		return nil, err
	}
	resp, err := cli.Do(req)
	if err != nil {
		logger.Errorf("request filrep.io api failed, %v", err)
		return nil, err
	}

	body := resp.Body
	defer body.Close()

	payload, err := ioutil.ReadAll(body)
	if err != nil {
		logger.Errorf("request filrep.io api failed, %v", err)
		return nil, err
	}

	var data MinerList
	if err = json.Unmarshal(payload, &data); err != nil {
		logger.WithField("payload", string(payload)).Warning("json unmarshal failed")
		logger.Errorf("request filrep.io api failed, %v", err)
		return nil, err
	}
	return &data, nil
}

func (c *Crawler) syncMinerList(ctx context.Context, limit, offset int) error {
	reqUrlFormat := c.filrepApi + "/miners?limit=%d&offset=%d"
	url := fmt.Sprintf(reqUrlFormat, limit, offset)
	data, err := c.minerList(ctx, url)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil && data.Pagination.Total < offset+limit {
			err = errors.New("finished")
		}
	}()
	if len(data.Miners) == 0 {
		return errors.New("finished")
	}

	list := make([]*models.Provider, 0, len(data.Miners))
	for _, miner := range data.Miners {
		item := &models.Provider{
			Address:         miner.Address,
			Status:          miner.Status,
			UptimeAverage:   miner.UptimeAverage,
			IsoCode:         miner.IsoCode,
			Region:          miner.Region,
			RawPower:        miner.RawPower,
			QualityAdjPower: miner.QualityAdjPower,
			FreeSpace:       miner.FreeSpace,

			Price:         miner.Price,
			VerifiedPrice: miner.VerifiedPrice,
			MinPieceSize:  miner.MinPieceSize,
			MaxPieceSize:  miner.MaxPieceSize,

			Score:                       miner.Score,
			ScoreUptime:                 miner.Scores.Uptime,
			ScoreStorageDeals:           miner.Scores.StorageDeals,
			ScoreCommittedSectorsProofs: miner.Scores.CommittedSectorsProofs,

			DealTotal:           miner.StorageDeals.Total,
			DealNoPenalties:     miner.StorageDeals.NoPenalties,
			DealAveragePrice:    miner.StorageDeals.AveragePrice,
			DealDataStored:      miner.StorageDeals.DataStored,
			DealSlashed:         miner.StorageDeals.Slashed,
			DealTerminated:      miner.StorageDeals.Terminated,
			DealFaultTerminated: miner.StorageDeals.FaultTerminated,
			DealSuccessRate:     miner.StorageDeals.SuccessRate,
		}

		list = append(list, item)
	}
	if err = c.repo.UpsertProvider(list); err != nil {
		logger.Errorf("UpsertProvider failed, %v", err)
	}
	return nil
}

func (c *Crawler) startDealTask(ctx context.Context) {
	logger.Info("sync deal data from s3...")
	defer logger.Info("sync deal data from s3 finished")

	url := "https://marketdeals.s3.amazonaws.com/StateMarketDeals.json"
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	req, err := c.makeRequest(ctx, url)
	if err != nil {
		logger.WithError(err).Errorf("makeRequest failed")
		return
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		logger.Errorf("request %s failed, %v", url, err)
		return
	}

	body := resp.Body
	defer body.Close()

	dec := json.NewDecoder(body)
	dealCh := make(chan *models.Deal, 1024)

	// skip '{'
	_, _ = dec.Token()
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			// handle error
			logger.Errorf("call Token failed, error:%v", err)
			break
		}

		var val filclient.MarketDeal
		err = dec.Decode(&val)
		if err != nil {
			// handle error
			logger.Errorf("call Decode failed, error:%v", err)
			break
		}

		dealId, err := strconv.ParseInt(fmt.Sprintf("%s", key), 10, 64)
		if err != nil {
			logger.WithField("key", key).Errorf("ParseInt error: %v", err)
			return
		}

		select {
		case <-ctx.Done():
			return
		case dealCh <- &models.Deal{
			DealId:               abi.DealID(dealId),
			PieceCid:             val.Proposal.PieceCID.String(),
			PieceSize:            val.Proposal.PieceSize,
			VerifiedDeal:         val.Proposal.VerifiedDeal,
			Client:               val.Proposal.Client.String(),
			Provider:             val.Proposal.Provider.String(),
			Label:                val.Proposal.Label,
			StartEpoch:           val.Proposal.StartEpoch,
			EndEpoch:             val.Proposal.EndEpoch,
			StoragePricePerEpoch: decimal.NewFromBigInt(val.Proposal.StoragePricePerEpoch.Int, 0),
			ProviderCollateral:   decimal.NewFromBigInt(val.Proposal.ProviderCollateral.Int, 0),
			ClientCollateral:     decimal.NewFromBigInt(val.Proposal.ClientCollateral.Int, 0),
			SectorStartEpoch:     val.State.SectorStartEpoch,
			LastUpdatedEpoch:     val.State.LastUpdatedEpoch,
			SlashEpoch:           val.State.SlashEpoch,
		}:
		default:
			deals := make([]*models.Deal, 0, 1024)
			empty := false
			for !empty {
				select {
				case deal := <-dealCh:
					deals = append(deals, deal)
				default:
					empty = true
					break
				}
			}
			err = c.repo.UpsertDeal(deals)
			if err != nil {
				logger.WithError(err).Error("call UpsertDeal failed")
				return
			}
			// Check whether any missing records need to be synchronized in the deal table
			//id := deals[len(deals)-1].DealId
			//if count, err := c.repo.CountDeal(uint64(id)); err != nil {
			//	logger.WithError(err).Error("call CountDeal failed")
			//	return
			//} else if int64(id)+1 == count {
			//	return
			//}
		}
	}
}
