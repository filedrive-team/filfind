package models

import (
	"github.com/filedrive-team/filfind/backend/types"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt types.UnixTime  `json:"created_at"`
	UpdatedAt types.UnixTime  `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

var autoMigrateModels = make([]interface{}, 0)

func Models() []interface{} {
	return autoMigrateModels
}
