package models

import (
	uuid "github.com/satori/go.uuid"
)

const (
	SPOwnerRole = "sp_owner"
	ClientRole  = "data_client"
	SystemRole  = "system"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &User{})
}

type User struct {
	Model
	Uid            uuid.UUID `json:"uid" gorm:"size:64;not null;uniqueIndex:uidx_user_uid"`
	Email          string    `json:"email" gorm:"size:255;not null;uniqueIndex:uidx_user_email"`
	HashedPassword string    `json:"password" gorm:"column:password;size:80"`
	Type           string    `json:"type" gorm:"size:32;not null;uniqueIndex:uidx_user_addr_type,priority:2"` // "sp_owner" or "data_client" or "system"
	AddressRobust  string    `json:"address_robust" gorm:"size:255;not null;uniqueIndex:uidx_user_addr_type,priority:1"`
	AddressId      string    `json:"address_id" gorm:"size:32"`
	Name           string    `json:"name" gorm:"size:128"`
	Avatar         string    `json:"avatar" gorm:"size:1024"`
	Logo           string    `json:"logo" gorm:"size:1024"`
	Location       string    `json:"location" gorm:"size:128"`
	ContactEmail   string    `json:"contact_email" gorm:"size:256"`
	Slack          string    `json:"slack" gorm:"size:128"`
	Github         string    `json:"github" gorm:"size:128"`
	Twitter        string    `json:"twitter" gorm:"size:128"`
	Description    string    `json:"description" gorm:"size:2048"`
}
