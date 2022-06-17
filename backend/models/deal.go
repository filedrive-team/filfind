package models

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/shopspring/decimal"
)

func init() {
	autoMigrateModels = append(autoMigrateModels, &Deal{})
}

type Deal struct {
	DealId       abi.DealID `gorm:"primarykey"`
	PieceCid     string     `gorm:"size:255"`
	PieceSize    abi.PaddedPieceSize
	VerifiedDeal bool
	Client       string `gorm:"size:255;index"`
	Provider     string `gorm:"size:32;index"`

	// Label is an arbitrary client chosen label to apply to the deal
	Label string `gorm:"size:255"`

	// Nominal start epoch. Deal payment is linear between StartEpoch and EndEpoch,
	// with total amount StoragePricePerEpoch * (EndEpoch - StartEpoch).
	// Storage deal must appear in a sealed (proven) sector no later than StartEpoch,
	// otherwise it is invalid.
	StartEpoch           abi.ChainEpoch `gorm:"index"`
	EndEpoch             abi.ChainEpoch
	StoragePricePerEpoch decimal.Decimal `gorm:"type:decimal(50,0)"`

	ProviderCollateral decimal.Decimal `gorm:"type:decimal(50,0)"`
	ClientCollateral   decimal.Decimal `gorm:"type:decimal(50,0)"`

	SectorStartEpoch abi.ChainEpoch // -1 if not yet included in proven sector
	LastUpdatedEpoch abi.ChainEpoch // -1 if deal state never updated
	SlashEpoch       abi.ChainEpoch // -1 if deal never slashed
}
