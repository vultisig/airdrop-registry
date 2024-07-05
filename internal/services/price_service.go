package services

import (
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func SavePrice(price *models.Price) error {
	return db.DB.Create(price).Error
}

func GetPriceByTokenAndDate(token string, date int64) (*models.Price, error) {
	var price models.Price
	err := db.DB.Where("token = ? AND date = ?", token, date).First(&price).Error
	return &price, err
}
