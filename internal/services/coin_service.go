package services

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/internal/models"
)

// AddCoin adds a coin to the vault
func (s *Storage) AddCoin(coin *models.CoinDBModel) error {
	if err := s.db.Create(coin).Error; err != nil {
		return fmt.Errorf("failed to add coin: %w", err)
	}
	return nil
}

// DeleteCoin deletes a coin by its ID , and the vault id
func (s *Storage) DeleteCoin(coinID string, vaultID uint) error {
	if err := s.db.Where("id = ? AND vault_id = ?", coinID, vaultID).Unscoped().Delete(&models.CoinDBModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete coin with ID %s: %w", coinID, err)
	}
	return nil
}

// GetCoin returns a coin by its ID
func (s *Storage) GetCoin(coinID string) (models.CoinDBModel, error) {
	var coin models.CoinDBModel
	if err := s.db.Where("id = ?", coinID).First(&coin).Error; err != nil {
		return coin, fmt.Errorf("failed to get coin with ID %s: %w", coinID, err)
	}
	return coin, nil
}

// GetCoins returns all coins for a vault
func (s *Storage) GetCoins(vaultID uint) ([]models.CoinDBModel, error) {
	var coins []models.CoinDBModel
	if err := s.db.Where("vault_id = ?", vaultID).Find(&coins).Error; err != nil {
		return coins, fmt.Errorf("failed to get coins for vault with ID %d: %w", vaultID, err)
	}
	return coins, nil
}
