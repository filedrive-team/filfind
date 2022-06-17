package models

import (
	"github.com/filedrive-team/filfind/backend/types"
	uuid "github.com/satori/go.uuid"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &LoginInfo{})
}

type LoginInfo struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      types.UnixTime `json:"created_at"`
	Uid            uuid.UUID      `json:"uid" gorm:"size:64;index:idx_logininfo_uid"`
	Ip             string         `json:"ip" gorm:"size:64"`
	Browser        string         `json:"browser" gorm:"size:32"`
	BrowserVersion string         `json:"browser_version" gorm:"size:64"`
	OS             string         `json:"os" gorm:"size:64"`
	Platform       string         `json:"platform" gorm:"size:32"`
	Mobile         bool           `json:"mobile"`
	Bot            bool           `json:"bot"`
	Renewal        bool           `json:"renewal"` // Renewal, not a user login action
}
