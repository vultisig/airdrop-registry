package models

import (
	"time"

	"gorm.io/gorm"
)

type Point struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ECDSA     string         `gorm:"type:varchar(255);index:idx_ecdsa_eddsa,priority:1" json:"ecdsa" binding:"required"`
	EDDSA     string         `gorm:"type:varchar(255);index:idx_ecdsa_eddsa,priority:2" json:"eddsa" binding:"required"`
	Balance   float64        `json:"balance"`
	Share     float64        `json:"share"`
	CycleID   uint           `json:"-"`
	Cycle     Cycle          `json:"cycle" gorm:"foreignKey:CycleID"`
}
