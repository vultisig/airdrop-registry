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

func GetPricesByToken(token string) ([]models.Price, error) {
	var prices []models.Price
	err := db.DB.Where("token = ?", token).Find(&prices).Error
	return prices, err
}

func GetPricesByTokenAndDateRange(token string, start int64, end int64) ([]models.Price, error) {
	var prices []models.Price
	err := db.DB.Where("token = ? AND date >= ? AND date <= ?", token, start, end).Find(&prices).Error
	return prices, err
}

func GetLatestPriceByToken(token string) (*models.Price, error) {
	var price models.Price
	err := db.DB.Where("token = ?", token).Order("date desc").First(&price).Error
	return &price, err
}

func GetLatestPrices() ([]models.Price, error) {
	var prices []models.Price
	err := db.DB.Raw("SELECT * FROM prices WHERE date IN (SELECT MAX(date) FROM prices GROUP BY token)").Scan(&prices).Error
	return prices, err
}
