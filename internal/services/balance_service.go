package services

import (
	"fmt"
	"time"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
	"gorm.io/gorm"
)

func SaveBalance(balance *models.Balance) error {
	if err := db.DB.Create(balance).Error; err != nil {
		return fmt.Errorf("failed to save balance: %w", err)
	}
	return nil
}

func SaveBalanceWithLatestPrice(balance *models.Balance) (float64, error) {
	if balance.Token == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		balance.Token = "ETH"
	}

	var latestPrice models.Price
	err := db.DB.Where("chain = ? AND token = ?", balance.Chain, balance.Token).
		Order("created_at DESC").
		First(&latestPrice).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("failed to fetch latest price: %w", err)
	}

	if latestPrice.ID != 0 && time.Since(latestPrice.CreatedAt) <= 24*time.Hour {
		balance.PriceID = latestPrice.ID
	}

	if err := db.DB.Create(balance).Error; err != nil {
		return 0, fmt.Errorf("failed to save balance with latest price: %w", err)
	}
	return latestPrice.Price, nil
}

func GetBalanceByID(id uint) (*models.Balance, error) {
	var balance models.Balance
	if err := db.DB.First(&balance, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get balance with id %d: %w", id, err)
	}
	return &balance, nil
}

func UpdateBalance(balance *models.Balance) error {
	if err := db.DB.Save(balance).Error; err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

func DeleteBalance(id uint) error {
	if err := db.DB.Delete(&models.Balance{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete balance with id %d: %w", id, err)
	}
	return nil
}

func GetAllBalances(page, pageSize int) ([]models.Balance, error) {
	var balances []models.Balance
	query := db.DB.Order("created_at DESC")

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&balances).Error; err != nil {
		return nil, fmt.Errorf("failed to get all balances: %w", err)
	}
	return balances, nil
}

func GetBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.BalanceResponse, error) {
	balances, err := fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances by vault keys: %w", err)
	}
	return balances, nil
}

func GetLatestBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.BalanceResponse, error) {
	balances, err := fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest balances by vault keys: %w", err)
	}
	return balances, nil
}

func GetUniqueActiveChainTokenPairs() ([]ChainTokenPair, error) {
	var pairs []ChainTokenPair
	query := `
		SELECT DISTINCT chain, token
		FROM balances
		WHERE id IN (
			SELECT MAX(id)
			FROM balances
			GROUP BY chain, token, ecdsa, eddsa
		) AND balance > 0
	`

	if err := db.DB.Raw(query).Scan(&pairs).Error; err != nil {
		return nil, fmt.Errorf("failed to get unique active chain token pairs: %w", err)
	}
	return pairs, nil
}

func GetAverageBalanceSince(ecdsaPublicKey, eddsaPublicKey string, since time.Time) (float64, error) {
	var totalBalance float64
	query := `
        SELECT SUM(avg_balance_usd) as total_balance
        FROM (
            SELECT
                b.chain,
                b.token,
                AVG(b.balance * COALESCE(p.price, (
                    SELECT price
                    FROM prices
                    WHERE chain = b.chain AND token = b.token AND created_at <= b.updated_at
                    ORDER BY created_at DESC
                    LIMIT 1
                ))) as avg_balance_usd
            FROM balances b
            LEFT JOIN prices p ON b.price_id = p.id
            WHERE b.ecdsa = ? AND b.eddsa = ? AND b.updated_at > ?
            GROUP BY b.chain, b.token
        ) as avg_balances
    `

	if err := db.DB.Raw(query, ecdsaPublicKey, eddsaPublicKey, since).Scan(&totalBalance).Error; err != nil {
		return 0, fmt.Errorf("failed to get average balance since %v: %w", since, err)
	}

	return totalBalance, nil
}

func GetAverageBalanceForTimeRange(ecdsaPublicKey, eddsaPublicKey string, startTime, endTime time.Time) (float64, error) {
	var totalBalance float64
	query := `
        SELECT SUM(avg_balance_usd) as total_balance
        FROM (
            SELECT
                b.chain,
                b.token,
                AVG(b.balance * COALESCE(p.price, (
                    SELECT price
                    FROM prices
                    WHERE chain = b.chain AND token = b.token AND created_at <= b.updated_at
                    ORDER BY created_at DESC
                    LIMIT 1
                ))) as avg_balance_usd
            FROM balances b
            LEFT JOIN prices p ON b.price_id = p.id
            WHERE b.ecdsa = ? AND b.eddsa = ? AND b.updated_at > ? AND b.updated_at <= ?
            GROUP BY b.chain, b.token
        ) as avg_balances
    `

	if err := db.DB.Raw(query, ecdsaPublicKey, eddsaPublicKey, startTime, endTime).Scan(&totalBalance).Error; err != nil {
		return 0, fmt.Errorf("failed to get average balance for time range %v to %v: %w", startTime, endTime, err)
	}

	return totalBalance, nil
}

type ChainTokenPair struct {
	Chain string
	Token string
}

func fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string, onlyRecent bool) ([]models.BalanceResponse, error) {
	type tempBalance struct {
		ID             uint
		ECDSA          string
		EDDSA          string
		Chain          string
		Address        string
		Token          string
		Balance        float64
		CreatedAt      time.Time
		UpdatedAt      time.Time
		UsdValue       float64
		PriceID        *uint
		PriceChain     *string
		PriceToken     *string
		PricePrice     *float64
		PriceCreatedAt *time.Time
		PriceUpdatedAt *time.Time
		PriceSource    *string
	}

	query := `
        SELECT
            b.id, b.ecdsa, b.eddsa, b.chain, b.address, b.token,
            b.balance, b.created_at, b.updated_at,
            COALESCE(b.balance * p.price, 0) as usd_value,
            p.id as price_id, p.chain as price_chain, p.token as price_token,
            p.price as price_price, p.created_at as price_created_at,
            p.updated_at as price_updated_at, p.source as price_source
        FROM balances b
        LEFT JOIN prices p ON b.price_id = p.id
        WHERE b.ecdsa = ? AND b.eddsa = ?
    `

	var args []interface{}
	args = append(args, ecdsaPublicKey, eddsaPublicKey)

	if onlyRecent {
		query += `
            AND b.id IN (
                SELECT MAX(id)
                FROM balances
                WHERE ecdsa = ? AND eddsa = ?
                GROUP BY chain, token
            )
            AND b.updated_at > ?
        `
		args = append(args, ecdsaPublicKey, eddsaPublicKey, time.Now().Add(-1*time.Hour))
	}

	query += " ORDER BY b.created_at DESC"

	var tempBalances []tempBalance
	if err := db.DB.Raw(query, args...).Scan(&tempBalances).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch balances by vault keys: %w", err)
	}

	var balances []models.BalanceResponse
	for _, tb := range tempBalances {
		br := models.BalanceResponse{
			ID:        tb.ID,
			CreatedAt: tb.CreatedAt,
			UpdatedAt: tb.UpdatedAt,
			ECDSA:     tb.ECDSA,
			EDDSA:     tb.EDDSA,
			Chain:     tb.Chain,
			Address:   tb.Address,
			Token:     tb.Token,
			Balance:   tb.Balance,
			UsdValue:  tb.UsdValue,
		}

		if tb.PriceID != nil {
			br.Price = &models.Price{
				ID:        *tb.PriceID,
				Chain:     *tb.PriceChain,
				Token:     *tb.PriceToken,
				Price:     *tb.PricePrice,
				CreatedAt: *tb.PriceCreatedAt,
				UpdatedAt: *tb.PriceUpdatedAt,
				Source:    *tb.PriceSource,
			}
		} else {
			recentPrice, err := GetLatestPriceByToken(tb.Chain, tb.Token)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("failed to get latest price for chain %s, token %s: %w", tb.Chain, tb.Token, err)
			}
			if recentPrice != nil {
				br.Price = recentPrice
				br.UsdValue = tb.Balance * recentPrice.Price
			}
		}

		balances = append(balances, br)
	}

	return balances, nil
}
