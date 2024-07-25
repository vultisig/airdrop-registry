package models

import (
	"time"
)

type Point struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ECDSA string `gorm:"type:varchar(255)" json:"ECDSA" binding:"required"`
	EDDSA string `gorm:"type:varchar(255)" json:"EdDSA" binding:"required"`

	Balance float64 `json:"balance"`
	Points  float64 `json:"points"`
}
