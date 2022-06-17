package models

import (
	"github.com/filedrive-team/filfind/backend/types"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &ProviderStatsMonthly{})
}

type ProviderStatsMonthly struct {
	ID       uint           `gorm:"primarykey" json:"id"`
	Month    types.UnixTime `json:"month" gorm:"uniqueIndex:idx_psm_sp_cli_m,priority:3"`
	Provider string         `json:"provider" gorm:"size:32;not null;uniqueIndex:idx_psm_sp_cli_m,priority:1"`
	Client   string         `json:"client" gorm:"size:32;not null;uniqueIndex:idx_psm_sp_cli_m,priority:2"`
}
