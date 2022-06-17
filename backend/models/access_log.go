package models

import (
	"github.com/filedrive-team/filfind/backend/types"
	"time"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &AccessLog{})
}

type AccessLog struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      types.UnixTime `json:"created_at" gorm:"index"`
	Method         string         `json:"method" gorm:"size:32"`
	Uri            string         `json:"uri" gorm:"size:1024;index"`
	FullPath       string         `json:"full_path" gorm:"size:1024"`
	Cost           time.Duration  `json:"cost"`
	Ip             string         `json:"ip" gorm:"size:64"`
	Browser        string         `json:"browser" gorm:"size:32"`
	BrowserVersion string         `json:"browser_version" gorm:"size:64"`
	OS             string         `json:"os" gorm:"size:64"`
	Platform       string         `json:"platform" gorm:"size:32"`
	Mobile         bool           `json:"mobile"`
	Bot            bool           `json:"bot"`
	Status         int            `json:"status"`
}
