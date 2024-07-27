package models

import (
	"time"
)

type Balance struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ECDSA   string  `gorm:"type:varchar(255)" json:"ecdsa" binding:"required"`
	EDDSA   string  `gorm:"type:varchar(255)" json:"eddsa" binding:"required"`
	Chain   string  `json:"chain"`
	Address string  `json:"address"`
	Token   string  `json:"token"`
	Balance float64 `json:"balance"`
	Date    int64   `json:"date"`
	PriceID uint    `json:"price_id"`
}

type BalanceResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ECDSA    string  `json:"ecdsa"`
	EDDSA    string  `json:"eddsa"`
	Chain    string  `json:"chain"`
	Address  string  `json:"address"`
	Token    string  `json:"token"`
	Balance  float64 `json:"balance"`
	Date     int64   `json:"date"`
	UsdValue float64 `json:"usd_value"`
	Price    *Price  `json:"price,omitempty"`
}
