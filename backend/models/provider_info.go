package models

import uuid "github.com/satori/go.uuid"

func init() {
	autoMigrateModels = append(autoMigrateModels, &ProviderInfo{})
}

type ProviderInfo struct {
	Model

	Uid             uuid.UUID `json:"uid" gorm:"size:64;not null;index"`
	Address         string    `json:"address" gorm:"size:32;not null;uniqueIndex:uidx_pinfo_addr"`
	AvailableDeals  string    `json:"available_deals" gorm:"size:128"`
	Bandwidth       string    `json:"bandwidth" gorm:"size:128"`
	SealingSpeed    string    `json:"sealing_speed" gorm:"size:128"`
	ParallelDeals   string    `json:"parallel_deals" gorm:"size:128"`
	RenewableEnergy string    `json:"renewable_energy" gorm:"size:128"`
	Certification   string    `json:"certification" gorm:"size:128"`
	IsMember        string    `json:"is_member" gorm:"size:16"`
	Experience      string    `json:"experience" gorm:"size:128"`
}
