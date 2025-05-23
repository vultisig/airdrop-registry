package models

import "gorm.io/gorm"

// Store vault rank and points for each season
type VaultSeasonStats struct {
	gorm.Model
	VaultID       uint    `gorm:"type:bigint;not null;uniqueIndex:vault_season_idx" json:"vault_id"`
	SeasonID      uint    `gorm:"type:bigint;not null;uniqueIndex:vault_season_idx" json:"season_id"`
	Rank          int64   `json:"rank"` // rank of the vault
	Points        float64 `json:"points"`
	Balance       int64   `gorm:"type:bigint;default:0" json:"balance"` // latest balance of the vault
	LPValue       int64   `gorm:"type:bigint;default:0" json:"lp_value"`
	SwapVolume    float64 `gorm:"type:bigint;default:0" json:"swap_volume"`
	NFTValue      int64   `gorm:"type:bigint;default:0" json:"nft_value"`
	ReferralCount int64   `gorm:"type:bigint;default:0" json:"referral_count"`
}

func (*VaultSeasonStats) TableName() string {
	return "vault_season_stats"
}
