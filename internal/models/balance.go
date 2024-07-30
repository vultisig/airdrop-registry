package models

import (
	"time"

	"gorm.io/gorm"
)

type Balance struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ECDSA     string         `gorm:"type:varchar(255);index:idx_ecdsa_eddsa,priority:1" json:"ecdsa" binding:"required"`
	EDDSA     string         `gorm:"type:varchar(255);index:idx_ecdsa_eddsa,priority:2" json:"eddsa" binding:"required"`
	Chain     string         `json:"chain"`
	Address   string         `json:"address"`
	Token     string         `json:"token"`
	Balance   float64        `json:"balance"`
	PriceID   uint           `json:"price_id"`
}

type BalanceResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
	ECDSA     string    `json:"ecdsa"`
	EDDSA     string    `json:"eddsa"`
	Chain     string    `json:"chain"`
	Address   string    `json:"address"`
	Token     string    `json:"token"`
	Balance   float64   `json:"balance"`
	UsdValue  float64   `json:"usd_value"`
	Price     *Price    `json:"price,omitempty"`
}
