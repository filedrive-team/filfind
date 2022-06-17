package models

import uuid "github.com/satori/go.uuid"

func init() {
	autoMigrateModels = append(autoMigrateModels, &UserRelation{})
}

type UserRelation struct {
	Model

	Uid     uuid.UUID `json:"uid" gorm:"size:64;uniqueIndex:uidx_userrel_uid_ptuid,priority:1"`
	Partner uuid.UUID `json:"partner" gorm:"size:64;uniqueIndex:uidx_userrel_uid_ptuid,priority:2"`
}
