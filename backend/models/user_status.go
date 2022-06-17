package models

import uuid "github.com/satori/go.uuid"

func init() {
	autoMigrateModels = append(autoMigrateModels, &UserStatus{})
}

type UserStatus struct {
	Model

	Uid    uuid.UUID `json:"uid" gorm:"size:64;uniqueIndex"`
	Online bool      `json:"online"`
}
