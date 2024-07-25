package models

import (
	"time"
)

type Point struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ECDSA string `gorm:"type:varchar(255)" json:"ecdsa" binding:"required"`
	EDDSA string `gorm:"type:varchar(255)" json:"eddsa" binding:"required"`

	Balance float64 `json:"balance"`
	Share   float64 `json:"share"`
}
