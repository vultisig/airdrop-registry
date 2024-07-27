package models

import (
	"time"
)

type Cycle struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Points []Point `json:"-" gorm:"foreignKey:CycleID"`
}
