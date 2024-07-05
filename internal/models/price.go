package models

import (
	"gorm.io/gorm"
)

type Price struct {
	gorm.Model
	Token string  `json:"token"`
	USD   float64 `json:"usd"`
	Date  int64   `json:"date"`
}
