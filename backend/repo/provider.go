package repo

import (
	"context"
	"fmt"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func (m *Manager) UpsertProvider(list []*models.Provider) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		UpdateAll: true,
	}).Omit("first_deal_time", "retrieval_success_rate", "owner", "review_score", "reviews").
		Create(list).Error
}

func (m *Manager) UpdateProviderOwner(list []*models.Provider) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"owner"}),
	}).Create(list).Error
}

type Provider struct {
	Address         string          `json:"address"`
	Region          string          `json:"region"`
	IsoCode         string          `json:"iso_code"`
	RawPower        decimal.Decimal `json:"raw_power"`
	QualityAdjPower decimal.Decimal `json:"quality_adj_power"`

	Price         *string          `json:"price"`
	VerifiedPrice *string          `json:"verified_price"`
	MinPieceSize  *decimal.Decimal `json:"min_piece_size"`
	MaxPieceSize  *decimal.Decimal `json:"max_piece_size"`

	StorageDeals         uint64  `json:"storage_deals" gorm:"column:deal_total"`
	StorageSuccessRate   float64 `json:"storage_success_rate" gorm:"column:deal_success_rate"`
	RetrievalSuccessRate float64 `json:"retrieval_success_rate"`

	FirstDealTime uint64 `json:"first_deal_time"`
	Owner         string `json:"owner"`

	Name            string  `json:"name"`
	ReputationScore int     `json:"reputation_score" gorm:"column:score"`
	ReviewScore     float64 `json:"review_score"`
	Reviews         int     `json:"reviews"`
	AvailableDeals  string  `json:"available_deals"`
	Bandwidth       string  `json:"bandwidth"`
	SealingSpeed    string  `json:"sealing_speed"`
	ParallelDeals   string  `json:"parallel_deals"`
	RenewableEnergy string  `json:"renewable_energy"`
	Certification   string  `json:"certification"`
	IsMember        string  `json:"is_member"`
	Experience      string  `json:"experience"`
}

type RangeParam struct {
	Min *decimal.Decimal `json:"min" form:"min" binding:"omitempty"`
	Max *decimal.Decimal `json:"max" form:"max" binding:"omitempty"`
}

type OrderParam struct {
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneof=reputation_score review_score storage_success_rate retrieval_success_rate price verified_price iso_code quality_adj_power storage_deals" example:"reputation_score"`
	Order  string `json:"order" form:"order" binding:"oneof=asc desc" example:"desc"`
}

type FilterParam struct {
	SpsStatus string `json:"sps_status" form:"sps_status" binding:"omitempty,oneof=all registered autofilled"`
	Region    string `json:"region" form:"region" binding:"omitempty,max=32"`

	RawPowerRange           *RangeParam `json:"raw_power_range" form:"raw_power_range" binding:"omitempty"`
	StorageSuccessRateRange *RangeParam `json:"storage_success_rate_range" form:"storage_success_rate_range" binding:"omitempty"`
	ReputationScoreRange    *RangeParam `json:"reputation_score_range" form:"reputation_score_range" binding:"omitempty"`
	ReviewScoreRange        *RangeParam `json:"review_score_range" form:"review_score_range" binding:"omitempty"`
}

func (m *Manager) ProviderList(pagination types.PaginationParams, order OrderParam, search string, filter FilterParam) (total int, list []*Provider, err error) {
	offset, size := utils.PaginationHelper(pagination.Page, pagination.PageSize, settings.DefaultPageSize)

	stmtCount := `
SELECT
count(*)
from provider p
left join user u 
on p.owner=u.address_id and u.type ='sp_owner'
where p.deleted_at is null
	`
	stmt := `
SELECT
p.*,
pi.available_deals,
pi.bandwidth,
pi.sealing_speed,
pi.parallel_deals,
pi.renewable_energy,
pi.certification,
pi.is_member,
pi.experience,
u.name
from provider p
left join provider_info pi
on p.address=pi.address
left join user u 
on p.owner=u.address_id and u.type ='sp_owner'
where p.deleted_at is null
	`
	var stmtParams []interface{}
	addCond := func(cond string, params ...interface{}) {
		stmtCount += cond
		stmt += cond
		stmtParams = append(stmtParams, params...)
	}
	if len(search) > 0 {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			logger.Fatal("validator engine changed")
			return
		}
		err = v.VarCtx(context.TODO(), search, "addressid")
		if err == nil {
			addCond(" and p.address = ?", search)
		} else {
			addCond(" and (p.iso_code like ? or name like ?)", "%"+search+"%", "%"+search+"%")
		}
	}
	if len(filter.SpsStatus) > 0 {
		var cond string
		switch filter.SpsStatus {
		case "registered":
			cond = " and name is not null"
		case "autofilled":
			cond = " and name is null"
		case "all":
		}
		if len(cond) > 0 {
			addCond(cond)
		}
	}
	if len(filter.Region) > 0 {
		if filter.Region != "all" {
			addCond(" and p.region=?", filter.Region)
		}
	}
	if filter.RawPowerRange != nil {
		if filter.RawPowerRange.Min != nil && !filter.RawPowerRange.Min.IsNegative() {
			addCond(" and p.raw_power>=?", *filter.RawPowerRange.Min)
		}
		if filter.RawPowerRange.Max != nil && !filter.RawPowerRange.Max.IsNegative() {
			addCond(" and p.raw_power<=?", *filter.RawPowerRange.Max)
		}
	}
	if filter.StorageSuccessRateRange != nil {
		if filter.StorageSuccessRateRange.Min != nil && !filter.StorageSuccessRateRange.Min.IsNegative() {
			addCond(" and p.deal_success_rate>=?", *filter.StorageSuccessRateRange.Min)
		}
		if filter.StorageSuccessRateRange.Max != nil && !filter.StorageSuccessRateRange.Max.IsNegative() {
			addCond(" and p.deal_success_rate<=?", *filter.StorageSuccessRateRange.Max)
		}
	}
	if filter.ReputationScoreRange != nil {
		if filter.ReputationScoreRange.Min != nil && !filter.ReputationScoreRange.Min.IsNegative() {
			addCond(" and p.score>=?", *filter.ReputationScoreRange.Min)
		}
		if filter.ReputationScoreRange.Max != nil && !filter.ReputationScoreRange.Max.IsNegative() {
			addCond(" and p.score<=?", *filter.ReputationScoreRange.Max)
		}
	}
	if filter.ReviewScoreRange != nil {
		if filter.ReviewScoreRange.Min != nil && !filter.ReviewScoreRange.Min.IsNegative() {
			addCond(" and p.review_score>=?", *filter.ReviewScoreRange.Min)
		}
		if filter.ReviewScoreRange.Max != nil && !filter.ReviewScoreRange.Max.IsNegative() {
			addCond(" and p.review_score<=?", *filter.ReviewScoreRange.Max)
		}
	}
	err = m.db.Raw(stmtCount, stmtParams...).Scan(&total).Error
	if err != nil {
		return total, nil, err
	}
	switch order.SortBy {
	case "reputation_score":
		order.SortBy = "score"
	case "storage_success_rate":
		order.SortBy = "deal_success_rate"
	case "storage_deals":
		order.SortBy = "deal_total"
	case "price", "verified_price":
		order.SortBy = order.SortBy + "+0"
	}
	stmt += fmt.Sprintf(" order by %s %s,name desc limit ? offset ?", order.SortBy, order.Order)
	stmtParams = append(stmtParams, size, offset)
	list = make([]*Provider, 0)
	err = m.db.Raw(stmt, stmtParams...).Scan(&list).Error
	return total, list, err
}

func (m *Manager) ProviderListByOwner(owner string) (list []*Provider, err error) {
	stmt := `
SELECT
p.*,
pi.available_deals,
pi.bandwidth,
pi.sealing_speed,
pi.parallel_deals,
pi.renewable_energy,
pi.certification,
pi.is_member,
pi.experience,
u.name
from provider p
left join provider_info pi
on p.address=pi.address
left join user u 
on p.owner=u.address_id and u.type ='sp_owner'
where p.deleted_at is null
and p.owner=?
`
	list = make([]*Provider, 0)
	err = m.db.Raw(stmt, owner).Scan(&list).Error
	return list, err
}

func (m *Manager) HasProvider(uid, provider string) (own bool, err error) {
	var count int64
	err = m.db.Model(models.User{}).
		Joins("left join provider p on user.address_id=p.owner").
		Where("p.address=? and user.uid=?", provider, uid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *Manager) GetProviderAddress(limit, offset int) (addresses []string, err error) {
	err = m.db.Model(models.Provider{}).
		Select("address").
		Limit(limit).
		Offset(offset).
		Scan(&addresses).Error
	return
}

func (m *Manager) UpdateProviderFirstDealTime(list []*models.Provider) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_deal_time"}),
	}).Create(list).Error
}

func (m *Manager) UpdateProviderReviewsStats(list []*models.Provider) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"review_score", "reviews"}),
	}).Create(list).Error
}

func (m *Manager) UpdateProviderSuccessRate(list []*models.Provider) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"deal_success_rate", "retrieval_success_rate"}),
	}).Create(list).Error
}

func (m *Manager) RegisteredSpRatio() (ratio float64, err error) {
	stmt := `
select COALESCE(sum(t.num),0)/count(*)
from provider p 
left join (
select address_id, 1 as num 
from user u 
where u.type ='sp_owner' 
and LENGTH(address_id)>0
) t on p.owner = t.address_id 
`
	err = m.db.Raw(stmt).Scan(&ratio).Error
	return
}
