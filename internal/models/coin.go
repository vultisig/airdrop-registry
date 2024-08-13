package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/vultisig/airdrop-registry/internal/common"
)

type CoinBase struct {
	Chain           common.Chain `json:"chain" binding:"required" gorm:"type:varchar(50);uniqueIndex:chain_ticker_address_idx;not null"`
	Ticker          string       `json:"ticker" binding:"required" gorm:"type:varchar(255);uniqueIndex:chain_ticker_address_idx;not null"`
	Address         string       `json:"address" binding:"required" gorm:"type:varchar(255);uniqueIndex:chain_ticker_address_idx;not null"`
	ContractAddress string       `json:"contract_address" gorm:"type:varchar(255)"`
	Decimals        int          `json:"decimals" binding:"required" gorm:"type:Integer;not null"`
	PriceProviderID string       `json:"price_provider_id" gorm:"type:varchar(255)"`
	IsNativeToken   bool         `json:"is_native_token" binding:"required"`
	HexPublicKey    string       `json:"hex_public_key" binding:"required" gorm:"type:varchar(255);not null"`
	Balance         string       `json:"balance" gorm:"type:varchar(50)"`
	Price           string       `json:"price" gorm:"type:varchar(50)"`
}

type CoinDBModel struct {
	ID        string `gorm:"type:varchar(255);not null;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CoinBase
	VaultID uint `json:"vault_id" binding:"required" gorm:"not null"`
}

func (CoinDBModel) TableName() string {
	return "coins"
}