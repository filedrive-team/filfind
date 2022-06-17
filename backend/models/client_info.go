package models

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &ClientInfo{})
}

type ClientInfo struct {
	Model

	Uid                uuid.UUID       `json:"uid" gorm:"size:64;not null;uniqueIndex:uidx_cinfo_uid"`
	AddressId          string          `json:"address_id" gorm:"size:32;not null;uniqueIndex:uidx_cinfo_addrid"`
	StorageCapacity    decimal.Decimal `json:"storage_capacity" gorm:"type:decimal(50,0)"`
	StorageDeals       int64           `json:"storage_deals"`
	DataCap            decimal.Decimal `json:"data_cap" gorm:"type:decimal(50,0)"` // available datacap
	UsedDataCap        decimal.Decimal `json:"used_data_cap" gorm:"type:decimal(50,0);default:0"`
	VerifiedDeals      int64           `json:"verified_deals"`
	Bandwidth          string          `json:"bandwidth" gorm:"size:128"`
	MonthlyStorage     string          `json:"monthly_storage" gorm:"size:128"`
	UseCase            string          `json:"use_case" gorm:"size:128"`
	ServiceRequirement string          `json:"service_requirement" gorm:"size:1024"`
}
