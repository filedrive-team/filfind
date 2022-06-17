package models

import "github.com/shopspring/decimal"

func init() {
	autoMigrateModels = append(autoMigrateModels, &Provider{})
}

type Provider struct {
	Model

	Address         string           `json:"address" gorm:"size:32;not null;uniqueIndex:uidx_provider_addr"`
	Status          bool             `json:"status"`         // true: reachable false: unreachable
	UptimeAverage   float64          `json:"uptime_average"` // reachable success rate
	IsoCode         string           `json:"iso_code" gorm:"size:4"`
	Region          string           `json:"region" gorm:"size:32"`
	RawPower        decimal.Decimal  `json:"raw_power" gorm:"type:decimal(50,0)"`
	QualityAdjPower decimal.Decimal  `json:"quality_adj_power" gorm:"type:decimal(50,0)"`
	FreeSpace       *decimal.Decimal `json:"free_space" gorm:"type:decimal(50,0)"`

	Price         *string          `json:"price" gorm:"size:80"`
	VerifiedPrice *string          `json:"verified_price" gorm:"size:80"`
	MinPieceSize  *decimal.Decimal `json:"min_piece_size" gorm:"type:decimal(50,0)"`
	MaxPieceSize  *decimal.Decimal `json:"max_piece_size" gorm:"type:decimal(50,0)"`

	Score                       decimal.Decimal `json:"score" gorm:"type:decimal(4,0)"`
	ScoreUptime                 decimal.Decimal `json:"score_uptime" gorm:"type:decimal(4,0)"`
	ScoreStorageDeals           decimal.Decimal `json:"score_storage_deals" gorm:"type:decimal(4,0)"`
	ScoreCommittedSectorsProofs decimal.Decimal `json:"score_committed_sectors_proofs" gorm:"type:decimal(4,0)"`

	DealTotal           uint64          `json:"deal_total"`
	DealNoPenalties     uint64          `json:"deal_no_penalties"`
	DealAveragePrice    decimal.Decimal `json:"deal_average_price" gorm:"type:decimal(50,0)"`
	DealDataStored      decimal.Decimal `json:"deal_data_stored" gorm:"type:decimal(50,0)"`
	DealSlashed         uint64          `json:"deal_slashed"`
	DealTerminated      uint64          `json:"deal_terminated"`
	DealFaultTerminated uint64          `json:"deal_fault_terminated"`
	DealSuccessRate     decimal.Decimal `json:"deal_success_rate" gorm:"type:decimal(16,12)"`

	FirstDealTime        int64           `json:"first_deal_time"`
	RetrievalSuccessRate decimal.Decimal `json:"retrieval_success_rate" gorm:"type:decimal(16,12)"`
	Owner                string          `json:"owner" gorm:"size:255;index"`
	ReviewScore          float64         `json:"review_score"`
	Reviews              uint64          `json:"reviews"`
}
