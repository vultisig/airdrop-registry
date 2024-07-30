package models

import (
	"time"

	"gorm.io/gorm"
)

type Cycle struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Points    []Point        `json:"-" gorm:"foreignKey:CycleID"`
}
