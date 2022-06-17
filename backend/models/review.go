package models

import (
	uuid "github.com/satori/go.uuid"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &Review{})
}

type Review struct {
	Model

	Uid      uuid.UUID `json:"uid" gorm:"size:64"`
	Client   string    `json:"client" gorm:"size:32"` // address id
	Provider string    `json:"provider" gorm:"size:32"`
	Score    int       `json:"score"`
	Content  string    `json:"content" gorm:"size:1024"`
	Title    string    `json:"title" gorm:"size:128"`
}
