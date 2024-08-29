package models

import (
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
	IsNativeToken   bool         `json:"is_native_token"`
	HexPublicKey    string       `json:"hex_public_key" binding:"required" gorm:"type:varchar(255);not null"`
	CMCId           int          `json:"cmc_id" gorm:"type:Integer"`
	Logo            string       `json:"logo" gorm:"type:varchar(255)"`
	Balance         string       `json:"balance" gorm:"type:varchar(50)"`
	PriceUSD        string       `json:"price" gorm:"type:varchar(50)"`
	USDValue        string       `json:"usd_value" gorm:"type:varchar(50)"`
}

type CoinDBModel struct {
	gorm.Model
	CoinBase
	VaultID uint `json:"vault_id" binding:"required" gorm:"not null"`
}

func (CoinDBModel) TableName() string {
	return "coins"
}

type CoinIdentity struct {
	Chain           common.Chain
	Ticker          string
	ContractAddress string
}
