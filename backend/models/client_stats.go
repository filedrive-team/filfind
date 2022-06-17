package models

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/shopspring/decimal"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &ClientStats{})
}

const SectionDeals int64 = 500000

type ClientStats struct {
	ID              uint            `gorm:"primarykey" json:"id"`
	CreatedAt       types.UnixTime  `json:"created_at"`
	UpdatedAt       types.UnixTime  `json:"updated_at"`
	SectionDealId   abi.DealID      `json:"section_deal_id" gorm:"not null;index;uniqueIndex:uidx_cs_cli_sec,priority:2"` // section statistical label
	Client          string          `json:"client" gorm:"size:32;not null;uniqueIndex:uidx_cs_cli_sec,priority:1"`
	UsedDataCap     decimal.Decimal `json:"used_data_cap" gorm:"type:decimal(50,0);default:0"`
	VerifiedDeals   int64           `json:"verified_deals" gorm:"default:0"`
	StorageCapacity decimal.Decimal `json:"storage_capacity" gorm:"type:decimal(50,0);default:0"`
	StorageDeals    int64           `json:"storage_deals" gorm:"default:0"`
}
