package services

import (
	"time"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func SaveBalance(balance *models.Balance) error {
	return db.DB.Create(balance).Error
}

func SaveBalanceWithLatestPrice(balance *models.Balance) (price float64, err error) {
	latestPrice, _ := GetLatestPriceByToken(balance.Chain, balance.Token)
	// if err != nil {
	// 	return 0, err
	// }

	if latestPrice.ID != 0 {
		if time.Now().Sub(latestPrice.CreatedAt) > 24*time.Hour {
			balance.PriceID = latestPrice.ID
		}
	}

	err = db.DB.Create(balance).Error
	return latestPrice.Price, err
}

func GetBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.Balance, error) {
	var balances []models.Balance
	err := db.DB.Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).Find(&balances).Error
	return balances, err
}

func GetLatestBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.BalanceResponse, error) {
	var balances []models.BalanceResponse

	oneHourAgo := time.Now().Add(-1 * time.Hour)

	subquery := db.DB.Table("balances").
		Select("MAX(id) as id").
		Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).
		Group("chain, token")

	err := db.DB.Table("balances").
		Select(`balances.id, balances.ecdsa, balances.eddsa, balances.chain, balances.address, balances.token,
				balances.balance, balances.date, balances.price_id, balances.created_at, balances.updated_at,
				balances.balance * prices.price as usd_balance,
				prices.chain as price_chain, prices.token as price_token, prices.price as price, prices.date as price_date, prices.source as price_source`).
		Joins("JOIN prices ON balances.price_id = prices.id").
		Where("balances.id IN (?)", subquery).
		Where("balances.updated_at > ?", oneHourAgo).
		Order("balances.created_at desc").
		Scan(&balances).Error

	if err != nil {
		return nil, err
	}

	// Map price fields to nested PriceResponse struct
	for i, balance := range balances {
		balances[i].Price = models.Price{
			Chain:  balance.Price.Chain,
			Token:  balance.Price.Token,
			Price:  balance.Price.Price,
			Date:   balance.Price.Date,
			Source: balance.Price.Source,
		}
	}

	return balances, nil
}

type ChainTokenPair struct {
	Chain string `json:"chain"`
	Token string `json:"token"`
}

func GetUniqueActiveChainTokenPairs() ([]ChainTokenPair, error) {
	var chainTokenPairs []ChainTokenPair

	subquery := db.DB.Table("balances").
		Select("MAX(id) as id").
		Group("chain, token, ecdsa, eddsa")

	err := db.DB.Table("balances").
		Select("DISTINCT chain, token").
		Where("id IN (?) AND balance > 0", subquery).
		Scan(&chainTokenPairs).Error

	return chainTokenPairs, err
}
