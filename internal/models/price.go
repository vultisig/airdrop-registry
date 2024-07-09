package models

import (
	"gorm.io/gorm"
)

type Price struct {
	gorm.Model
	Chain  string  `json:"chain"`
	Token  string  `json:"token"`
	Price  float64 `json:"usd"`
	Date   int64   `json:"date"`
	Source string  `json:"source"`
}
