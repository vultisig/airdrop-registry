package services

import (
	"fmt"
	"time"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
	"gorm.io/gorm"
)

func SavePrice(price *models.Price) error {
	if err := db.DB.Create(price).Error; err != nil {
		return fmt.Errorf("failed to save price: %w", err)
	}
	return nil
}

func GetPriceByID(id uint) (*models.Price, error) {
	var price models.Price
	if err := db.DB.First(&price, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get price with id %d: %w", id, err)
	}
	return &price, nil
}

func UpdatePrice(price *models.Price) error {
	if err := db.DB.Save(price).Error; err != nil {
		return fmt.Errorf("failed to update price: %w", err)
	}
	return nil
}

func DeletePrice(id uint) error {
	if err := db.DB.Delete(&models.Price{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete price with id %d: %w", id, err)
	}
	return nil
}

func GetAllPrices(page, pageSize int) ([]models.Price, error) {
	var prices []models.Price
	query := db.DB.Order("created_at DESC")
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	if err := query.Find(&prices).Error; err != nil {
		return nil, fmt.Errorf("failed to get all prices: %w", err)
	}
	return prices, nil
}

func GetPriceByTokenAndTime(chain, token string, timestamp time.Time) (*models.Price, error) {
	var price models.Price
	err := db.DB.Where("chain = ? AND token = ? AND created_at = ?", chain, token, timestamp).First(&price).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get price for chain %s, token %s, time %v: %w", chain, token, timestamp, err)
	}
	return &price, nil
}

func GetPricesByToken(chain, token string) ([]models.Price, error) {
	var prices []models.Price
	err := db.DB.Where("chain = ? AND token = ?", chain, token).Find(&prices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get prices for chain %s, token %s: %w", chain, token, err)
	}
	return prices, nil
}

func GetPricesByTokenAndTimeRange(chain, token string, start, end time.Time) ([]models.Price, error) {
	var prices []models.Price
	err := db.DB.Where("chain = ? AND token = ? AND created_at >= ? AND created_at <= ?", chain, token, start, end).Find(&prices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get prices for chain %s, token %s, time range %v to %v: %w", chain, token, start, end, err)
	}
	return prices, nil
}

func GetLatestPriceByToken(chain, token string) (*models.Price, error) {
	var price models.Price
	err := db.DB.Where("chain = ? AND token = ?", chain, token).
		Order("created_at DESC").
		First(&price).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest price for chain %s, token %s: %w", chain, token, err)
	}
	return &price, nil
}

func GetLatestPrices() ([]models.Price, error) {
	var prices []models.Price
	err := db.DB.Raw("SELECT * FROM prices WHERE (chain, token, created_at) IN (SELECT chain, token, MAX(created_at) FROM prices GROUP BY chain, token)").Scan(&prices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get latest prices: %w", err)
	}
	return prices, nil
}
