package models

import (
	"fmt"
	"strings"

	"github.com/vultisig/airdrop-registry/pkg/utils"
	"gorm.io/gorm"
)

type Vault struct {
	ECDSA        string `gorm:"type:varchar(255);primaryKey" json:"ECDSA" binding:"required"`
	EDDSA        string `gorm:"type:varchar(255);primaryKey" json:"EdDSA" binding:"required"`
	HexChainCode string `gorm:"type:varchar(255)" json:"hexChainCode" binding:"required"`
}

func (Vault) TableName() string {
	return "vaults"
}

func (v *Vault) BeforeCreate(tx *gorm.DB) (err error) {
	v.ECDSA = strings.ToLower(v.ECDSA)
	v.EDDSA = strings.ToLower(v.EDDSA)

	if !utils.IsValidHex(v.ECDSA) || !utils.IsValidHex(v.EDDSA) {
		return fmt.Errorf("invalid ECDSA or EdDSA format")
	}

	if !utils.IsValidHex(v.HexChainCode) {
		return fmt.Errorf("invalid hex chain code format")
	}

	var count int64
	tx.Model(&Vault{}).Where("ecdsa = ? AND eddsa = ?", v.ECDSA, v.EDDSA).Count(&count)
	if count > 0 {
		return fmt.Errorf("vault with given ECDSA and EDDSA already exists")
	}
	return
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Vault{})
}
