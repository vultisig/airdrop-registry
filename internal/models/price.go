package models

import (
	"time"

	"gorm.io/gorm"
)

type Price struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Chain     string         `gorm:"index:idx_chain_token_date,priority:1" json:"chain"`
	Token     string         `gorm:"index:idx_chain_token_date,priority:2" json:"token"`
	Price     float64        `json:"usd"`
	Source    string         `json:"source"`
}
