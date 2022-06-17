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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (m *Manager) UpsertClientInfo(info *models.ClientInfo) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		UpdateAll: true,
	}).Omit("data_cap", "used_data_cap").
		Create(info).Error
}

type ClientInfo struct {
	AddressId          string          `json:"address_id"`
	StorageCapacity    decimal.Decimal `json:"storage_capacity"`
	StorageDeals       int64           `json:"storage_deals"`
	DataCap            decimal.Decimal `json:"data_cap"`
	UsedDataCap        decimal.Decimal `json:"used_data_cap"`
	TotalDataCap       decimal.Decimal `json:"total_data_cap"`
	VerifiedDeals      int64           `json:"verified_deals"`
	Bandwidth          string          `json:"bandwidth"`
	MonthlyStorage     string          `json:"monthly_storage"`
	UseCase            string          `json:"use_case"`
	ServiceRequirement string          `json:"service_requirement"`
}

func (m *Manager) ClientInfo(addrId string) (info *ClientInfo, err error) {
	info = &ClientInfo{}
	err = m.db.Model(&models.ClientInfo{}).Where("address_id=?", addrId).
		Select("*,(data_cap+used_data_cap) as total_data_cap").
		Scan(info).Error
	if err == nil && info.AddressId == "" {
		err = gorm.ErrRecordNotFound
	}
	return info, err
}

func (m *Manager) UpsertClientInfoDataCap(list []*models.ClientInfo) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"uid", "data_cap"}),
	}).Create(list).Error
}

func (m *Manager) UpsertClientInfoDeals(list []*models.ClientInfo) error {
	return m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"storage_capacity", "storage_deals", "used_data_cap", "verified_deals"}),
	}).Create(list).Error
}

type Client struct {
	ClientInfo
	Name          string `json:"name"`
	Location      string `json:"location"`
	AddressRobust string `json:"address_robust"`
}

type ClientOrderParam struct {
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneof=storage_deals storage_capacity total_data_cap used_data_cap data_cap verified_deals" example:"data_cap"`
	Order  string `json:"order" form:"order" binding:"oneof=asc desc" example:"desc"`
}

func (m *Manager) ClientList(pagination types.PaginationParams, order ClientOrderParam, search string) (total int, list []*Client, err error) {
	offset, size := utils.PaginationHelper(pagination.Page, pagination.PageSize, settings.DefaultPageSize)

	stmtCount := `
SELECT
count(*)
from client_info ci
left join user u 
on ci.address_id=u.address_id and u.type ='data_client'
where ci.deleted_at is null
	`
	stmt := `
select * from (
	SELECT
	ci.*,
	(ci.data_cap+ci.used_data_cap) as total_data_cap,
	u.name,
	u.location,
	u.address_robust
	from client_info ci
	left join user u 
	on ci.address_id=u.address_id and u.type ='data_client'
	where ci.deleted_at is null
) ci 
where 1
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
		if err = v.VarCtx(context.TODO(), search, "addressid"); err == nil {
			addCond(" and ci.address_id = ?", search)
		} else if err = v.VarCtx(context.TODO(), search, "address"); err == nil {
			addCond(" and address_robust = ?", search)
		} else {
			addCond(" and (location like ? or name like ?)", "%"+search+"%", "%"+search+"%")
		}
	}

	err = m.db.Raw(stmtCount, stmtParams...).Scan(&total).Error
	if err != nil {
		return total, nil, err
	}

	stmt += fmt.Sprintf(" order by %s %s limit ? offset ?", order.SortBy, order.Order)
	stmtParams = append(stmtParams, size, offset)
	list = make([]*Client, 0)
	err = m.db.Raw(stmt, stmtParams...).Scan(&list).Error
	return total, list, err
}
