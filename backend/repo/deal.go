package repo

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
	"time"
)

func (m *Manager) UpsertDeal(deals []*models.Deal) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "deal_id"}},
		UpdateAll: false,
	}).Create(deals).Error
}

// deal id range [0, lastDealId]
func (m *Manager) CountDeal(lastDealId uint64) (count int64, err error) {
	return m.CountDealRange(0, lastDealId+1)
}

// deal id range [firstDealId, lastDealId)
func (m *Manager) CountDealRange(firstDealId, lastDealId uint64) (count int64, err error) {
	err = m.db.Model(&models.Deal{}).
		Where("deal_id>=? and deal_id<?", firstDealId, lastDealId).
		Count(&count).Error
	return
}

// deal id range [firstDealId, lastDealId)
func (m *Manager) GetDealIdByRange(firstDealId, lastDealId uint64) (dealIds []uint64, err error) {
	err = m.db.Model(&models.Deal{}).
		Where("deal_id>=? and deal_id<?", firstDealId, lastDealId).
		Select("deal_id").
		Order("deal_id").
		Scan(&dealIds).Error
	return
}

func (m *Manager) GetDeal(dealId uint64) (deal *models.Deal, err error) {
	deal = &models.Deal{}
	err = m.db.Model(&models.Deal{}).Where("deal_id=?", dealId).First(deal).Error
	return
}

func (m *Manager) FirstDeal() (deal *models.Deal, err error) {
	deal = &models.Deal{}
	err = m.db.Model(&models.Deal{}).Order("deal_id").First(deal).Error
	return
}

func (m *Manager) FirstValidDeal() (deal *models.Deal, err error) {
	deal = &models.Deal{}
	err = m.db.Model(&models.Deal{}).Where("sector_start_epoch>0").Order("deal_id").First(deal).Error
	return
}

func (m *Manager) LastDeal() (deal *models.Deal, err error) {
	deal = &models.Deal{}
	err = m.db.Model(&models.Deal{}).Order("deal_id desc").First(deal).Error
	return
}

type ProviderDealFirstTime struct {
	Provider   string
	StartEpoch abi.ChainEpoch
}

func (m *Manager) StatsProviderFirstDealTime() (data []*ProviderDealFirstTime, err error) {
	stmt := `
select * from (
	SELECT provider,
	min(start_epoch) as start_epoch 
	FROM (
		select p.address from provider p 
		where p.first_deal_time is null
	) t left join deal d on t.address=d.provider and sector_start_epoch>0
	GROUP BY provider
) tt
where tt.provider is not null
`
	err = m.db.Raw(stmt).Scan(&data).Error
	return
}

type ClientHistoryDealStatsItem struct {
	Provider        string          `json:"provider"`
	Region          string          `json:"region"`
	IsoCode         string          `json:"iso_code"`
	RawPower        decimal.Decimal `json:"raw_power"`
	QualityAdjPower decimal.Decimal `json:"quality_adj_power"`

	MinPieceSize *decimal.Decimal `json:"min_piece_size"`
	MaxPieceSize *decimal.Decimal `json:"max_piece_size"`

	RetrievalSuccessRate float64 `json:"retrieval_success_rate"`

	Owner string `json:"owner"`

	StorageSuccessRate float64 `json:"storage_success_rate" gorm:"column:deal_success_rate"`
	ReputationScore    int     `json:"reputation_score" gorm:"column:score"`
	ReviewScore        float64 `json:"review_score"`
	Reviews            int     `json:"reviews"`

	Name string `json:"name"`
	//AvailableDeals     string  `json:"available_deals"`
	//Bandwidth          string  `json:"bandwidth"`
	//SealingSpeed       string  `json:"sealing_speed"`
	//ParallelDeals      string  `json:"parallel_deals"`
	//RenewableEnergy    string  `json:"renewable_energy"`
	//Certification      string  `json:"certification"`
	//IsMember           string  `json:"is_member"`
	//Experience         string  `json:"experience"`

	// relevant client deal
	StorageDeals     uint64           `json:"storage_deals"`
	StorageCapacity  decimal.Decimal  `json:"storage_capacity"`
	DataCap          decimal.Decimal  `json:"data_cap"`
	AvgPrice         *decimal.Decimal `json:"avg_price"`
	AvgVerifiedPrice *decimal.Decimal `json:"avg_verified_price"`
	FirstDealTime    uint64           `json:"first_deal_time" gorm:"column:client_first_deal_time"`
}

func (m *Manager) StatsClientHistoryDeal(pagination types.PaginationParams, clientId string) (total int, list []*ClientHistoryDealStatsItem, err error) {
	offset, size := utils.PaginationHelper(pagination.Page, pagination.PageSize, settings.DefaultPageSize)

	stmtCount := `
select 
count(*) 
from (
select provider from deal where client =? group by provider
) t
	`
	stmt := `
select 
t1.*,
t2.data_cap,
t2.avg_verified_price,
t3.avg_price,
p.*,
u.name
from (
select provider,count(*) as storage_deals,sum(piece_size) as storage_capacity, 
-- epoch change to unix time
1598306400+30*min(start_epoch) as client_first_deal_time
from deal where client=? group by provider
) t1 left join (
select provider,sum(piece_size) as data_cap, avg(storage_price_per_epoch) as avg_verified_price
from deal where verified_deal and client=? group by provider
) t2 on t1.provider=t2.provider
left join (
select provider,avg(storage_price_per_epoch) as avg_price
from deal where not verified_deal and client=? group by provider
) t3 on t1.provider=t3.provider
left join provider p
on t1.provider=p.address
left join user u
on p.owner=u.address_id and u.type=?
order by t1.storage_capacity desc
limit ? offset ?
	`

	err = m.db.Raw(stmtCount, clientId).Scan(&total).Error
	if err != nil {
		return total, nil, err
	}

	err = m.db.Raw(stmt, clientId, clientId, clientId, models.SPOwnerRole, size, offset).Scan(&list).Error
	return total, list, err
}

func (m *Manager) HasDeal(client, provider string) (has bool, err error) {
	var count int64
	err = m.db.Model(models.Deal{}).
		Where("client=? and provider=?", client, provider).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// deal id range [startDealId,endDealId)
func (m *Manager) StatsClientDeal(startDealId, endDealId int64) (list []*models.ClientStats, err error) {
	stmt := `
select t1.client,
t1.storage_capacity,
t1.storage_deals,
COALESCE(t2.used_data_cap,0) as used_data_cap,
COALESCE(t2.verified_deals,0) as verified_deals
from (
	select client, 
	sum(piece_size) as storage_capacity, 
	count(*) as storage_deals 
	from deal
	where deal_id>=@startDealId and deal_id<@endDealId
	group by client
) t1 left join (
	select client, 
	sum(piece_size) as used_data_cap, 
	count(*) as verified_deals 
	from deal
	where verified_deal and deal_id>=@startDealId and deal_id<@endDealId
	group by client
) t2 on t1.client=t2.client
`
	err = m.db.Raw(stmt, map[string]interface{}{
		"startDealId": startDealId,
		"endDealId":   endDealId,
	}).Scan(&list).Error
	return
}

func (m *Manager) ClientsToDealSpLastMonth() (count int64, err error) {
	begin := utils.MonthBegin(time.Now())
	end := begin.AddDate(0, 1, 0)
	return m.ClientsToDealSpByRange(utils.GetEpochByTime(begin), utils.GetEpochByTime(end))
}

func (m *Manager) ClientsToDealSpByRange(startEpoch, endEpoch int64) (count int64, err error) {
	stmt := `
select count(distinct client) from deal 
where start_epoch >= ? and start_epoch < ?
and client in (select ci.address_id from client_info ci where ci.data_cap+ci.used_data_cap>=100*pow(1024,4))
`
	err = m.db.Raw(stmt, startEpoch, endEpoch).Scan(&count).Error
	return
}

type ProviderClient struct {
	Provider string `json:"provider"`
	Client   string `json:"client"`
}

func (m *Manager) ProviderClientsByRange(startEpoch, endEpoch int64) (list []*ProviderClient, err error) {
	err = m.db.Model(&models.Deal{}).
		Where("start_epoch >=? and start_epoch <?", startEpoch, endEpoch).
		Group("provider,client").
		Select("provider,client").
		Scan(&list).Error
	return
}
