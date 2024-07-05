package models

import (
	"gorm.io/gorm"
)

type Balance struct {
	gorm.Model
	VaultID uint    `json:"vault_id"`
	Chain   string  `json:"chain"`
	Token   string  `json:"token"`
	Balance float64 `json:"balance"`
	Date    int64   `json:"date"`
	PriceID uint    `json:"price_id"`
}
