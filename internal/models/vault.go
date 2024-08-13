package models

import (
	"errors"

	"gorm.io/gorm"
)

var ErrAlreadyExist = errors.New("already exist")

type Vault struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(255)" json:"name" binding:"required"`
	ECDSA        string  `gorm:"type:varchar(255);uniqueIndex:ecdsa_eddsa_idx;not null" json:"ecdsa" binding:"required"`
	EDDSA        string  `gorm:"type:varchar(255);uniqueIndex:ecdsa_eddsa_idx;not null" json:"eddsa" binding:"required"`
	HexChainCode string  `gorm:"type:varchar(255)" json:"hex_chain_code" binding:"required"`
	Uid          string  `gorm:"type:varchar(255)" json:"uid" binding:"required"`
	TotalPoint   float64 `json:"total_point"`  // total point of the vault
	JoinAirdrop  bool    `json:"join_airdrop"` // join airdrop or not
}

func (*Vault) TableName() string {
	return "vaults"
}
