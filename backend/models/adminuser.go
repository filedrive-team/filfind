package models

import (
	uuid "github.com/satori/go.uuid"
)

const (
	AdminRole = "admin"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &AdminUser{})
}

type AdminUser struct {
	Model
	Uid            uuid.UUID `json:"uid" gorm:"size:64;not null;uniqueIndex:uidx_adminuser_uid"`
	Name           string    `json:"name" gorm:"size:128;not null;uniqueIndex:uidx_adminuser_name"`
	HashedPassword string    `json:"password" gorm:"column:password;size:80"`
	Type           string    `json:"type" gorm:"size:32"` // "admin"
	Avatar         string    `json:"avatar" gorm:"size:1024"`
	Description    string    `json:"description" gorm:"size:2048"`
}
