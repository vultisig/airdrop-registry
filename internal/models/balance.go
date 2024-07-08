package models

import (
	"gorm.io/gorm"
)

type Balance struct {
	gorm.Model
	ECDSA   string  `gorm:"type:varchar(255)" json:"ECDSA" binding:"required"`
	EDDSA   string  `gorm:"type:varchar(255)" json:"EdDSA" binding:"required"`
	Chain   string  `json:"chain"`
	Address string  `json:"address"`
	Token   string  `json:"token"`
	Balance float64 `json:"balance"`
	Date    int64   `json:"date"`
	PriceID uint    `json:"price_id"`
}
