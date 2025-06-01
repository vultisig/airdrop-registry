package services

import (
	"context"
	"fmt"
	"time"

	"github.com/vultisig/airdrop-registry/internal/common"
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

// DeleteCoin deletes a coin by its ID , and the vault id
func (s *Storage) DeleteCoins(coinIDs []uint, vaultID uint) error {
	if err := s.db.Where("id IN ? AND vault_id = ?", coinIDs, vaultID).Unscoped().Delete(&models.CoinDBModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete coins with IDs %v: %w", coinIDs, err)
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
func (s *Storage) UpdateCoinPrice(chain common.Chain, ticker string, priceUSD float64) error {
	qry := `UPDATE coins SET price_usd = ? WHERE chain = ? AND ticker = ?`
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := s.db.WithContext(ctx).Exec(qry, priceUSD, chain.String(), ticker).Error; err != nil {
		return fmt.Errorf("failed to update coin price: %w", err)
	}
	return nil
}

func (s *Storage) UpdateCoinPriceByCMCID(cmcID int, priceUSD float64) error {
	qry := `UPDATE coins SET price_usd = ? WHERE cmc_id = ? `
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := s.db.WithContext(ctx).Exec(qry, priceUSD, cmcID).Error; err != nil {
		return fmt.Errorf("failed to update coin price: %w", err)
	}
	return nil
}

func (s *Storage) GetUniqueCoins() ([]models.CoinIdentity, error) {
	var coinIdentities []models.CoinIdentity
	// if the query takes more than 2 minutes, it will be cancelled
	// let's investigate why it take so long to finish
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	qry := "SELECT DISTINCT chain, ticker, contract_address,cmc_id FROM coins where vault_id in (select id from vaults where join_airdrop = 1)"
	if err := s.db.WithContext(ctx).Raw(qry).Scan(&coinIdentities).Error; err != nil {
		return nil, fmt.Errorf("failed to get unique coins: %w", err)
	}
	return coinIdentities, nil
}

func (s *Storage) GetCoinsWithPage(startId, limit uint64) ([]models.CoinDBModel, error) {
	var coins []models.CoinDBModel
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := s.db.WithContext(ctx).Model(&models.CoinDBModel{}).Where("id > ?", startId).Limit(int(limit)).Scan(&coins).Error; err != nil {
		return coins, fmt.Errorf("failed to get coins: %w", err)
	}
	return coins, nil
}
func (s *Storage) UpdateCoinBalance(coinID uint64, balance float64) error {
	qry := `UPDATE coins SET balance = ?, usd_value = balance * price_usd  WHERE id = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.db.WithContext(ctx).Exec(qry, balance, coinID).Error; err != nil {
		return fmt.Errorf("failed to update coin balance: %w", err)
	}
	return nil
}
