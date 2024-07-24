package services

import (
	"strings"
	"time"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func SaveBalance(balance *models.Balance) error {
	return db.DB.Create(balance).Error
}

func SaveBalanceWithLatestPrice(balance *models.Balance) (float64, error) {
	if balance.Token == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		balance.Token = "ETH"
	}

	latestPrice, err := GetLatestPriceByToken(balance.Chain, balance.Token)
	if err != nil {
		recordNotFound := strings.Contains(err.Error(), "record not found")
		if !recordNotFound {
			return 0, err
		}
	}

	if latestPrice.ID != 0 {
		twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

		if latestPrice.CreatedAt.After(twentyFourHoursAgo) {
			balance.PriceID = latestPrice.ID
		}
	}

	err = db.DB.Create(balance).Error
	return latestPrice.Price, err
}

func GetBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.BalanceResponse, error) {
	return fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey, false)
}

func GetLatestBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.BalanceResponse, error) {
	return fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey, true)
}

func fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string, onlyRecent bool) ([]models.BalanceResponse, error) {
	var balances []models.BalanceResponse

	type balanceWithPrice struct {
		ID             uint       `json:"id"`
		ECDSA          string     `json:"ecdsa"`
		EDDSA          string     `json:"eddsa"`
		Chain          string     `json:"chain"`
		Address        string     `json:"address"`
		Token          string     `json:"token"`
		Balance        float64    `json:"balance"`
		Date           int64      `json:"date"`
		PriceID        *uint      `json:"price_id"`
		UsdBalance     float64    `json:"usd_balance"`
		CreatedAt      time.Time  `json:"created_at"`
		UpdatedAt      time.Time  `json:"updated_at"`
		PriceChain     *string    `json:"price_chain"`
		PriceToken     *string    `json:"price_token"`
		PricePrice     *float64   `json:"price_price"`
		PriceDate      *int64     `json:"price_date"`
		PriceSource    *string    `json:"price_source"`
		PriceCreatedAt *time.Time `json:"price_created_at"`
		PriceUpdatedAt *time.Time `json:"price_updated_at"`
	}

	var results []balanceWithPrice

	query := db.DB.Table("balances").
		Select(`balances.id, balances.ecdsa, balances.eddsa, balances.chain, balances.address, balances.token,
				balances.balance, balances.date, balances.price_id, balances.created_at, balances.updated_at,
				COALESCE(balances.balance * prices.price, 0) as usd_balance,
				prices.id as price_id, prices.chain as price_chain, prices.token as price_token, prices.price as price_price, prices.date as price_date, prices.source as price_source, prices.created_at as price_created_at, prices.updated_at as price_updated_at`).
		Joins("LEFT JOIN prices ON balances.price_id = prices.id").
		Where("balances.ecdsa = ? AND balances.eddsa = ?", ecdsaPublicKey, eddsaPublicKey)

	if onlyRecent {
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		subquery := db.DB.Table("balances").
			Select("MAX(id) as id").
			Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).
			Group("chain, token")

		query = query.Where("balances.id IN (?)", subquery).
			Where("balances.updated_at > ?", oneHourAgo)
	}

	err := query.Order("balances.created_at desc").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		balance := models.BalanceResponse{
			ID:         result.ID,
			ECDSA:      result.ECDSA,
			EDDSA:      result.EDDSA,
			Chain:      result.Chain,
			Address:    result.Address,
			Token:      result.Token,
			Balance:    result.Balance,
			Date:       result.Date,
			UsdBalance: result.UsdBalance,
			CreatedAt:  result.CreatedAt,
			UpdatedAt:  result.UpdatedAt,
			Price: models.PriceResponse{
				ID:        0,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				Chain:     "",
				Token:     "",
				Price:     0,
				Date:      0,
				Source:    "",
			},
		}
		if result.PriceID != nil {
			balance.Price.ID = *result.PriceID
			if result.PriceCreatedAt != nil {
				balance.Price.CreatedAt = *result.PriceCreatedAt
			}
			if result.PriceUpdatedAt != nil {
				balance.Price.UpdatedAt = *result.PriceUpdatedAt
			}
			if result.PriceChain != nil {
				balance.Price.Chain = *result.PriceChain
			}
			if result.PriceToken != nil {
				balance.Price.Token = *result.PriceToken
			}
			if result.PricePrice != nil {
				balance.Price.Price = *result.PricePrice
			}
			if result.PriceDate != nil {
				balance.Price.Date = *result.PriceDate
			}
			if result.PriceSource != nil {
				balance.Price.Source = *result.PriceSource
			}
		} else {
			recentPrice, err := GetLatestPriceByToken(balance.Chain, balance.Token)
			if err != nil {
				recordNotFound := strings.Contains(err.Error(), "record not found")
				if !recordNotFound {
					return nil, err
				}
			}
			if recentPrice != nil {
				balance.Price.ID = recentPrice.ID
				balance.Price.CreatedAt = recentPrice.CreatedAt
				balance.Price.UpdatedAt = recentPrice.UpdatedAt
				balance.Price.Chain = recentPrice.Chain
				balance.Price.Token = recentPrice.Token
				balance.Price.Price = recentPrice.Price
				balance.Price.Date = recentPrice.Date
				balance.Price.Source = recentPrice.Source

				balance.UsdBalance = balance.Balance * recentPrice.Price
			}
		}
		balances = append(balances, balance)
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
