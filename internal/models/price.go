package models

import (
	"time"
)

type Price struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Chain  string  `json:"chain"`
	Token  string  `json:"token"`
	Price  float64 `json:"usd"`
	Date   int64   `json:"date"`
	Source string  `json:"source"`
}
