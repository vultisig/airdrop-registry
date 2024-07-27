package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
	"gorm.io/gorm"
)

func SaveBalance(balance *models.Balance) error {
	return db.DB.Create(balance).Error
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
		return 0, err
	}

	if latestPrice.ID != 0 && time.Since(latestPrice.CreatedAt) <= 24*time.Hour {
		balance.PriceID = latestPrice.ID
	}

	err = db.DB.Create(balance).Error
	return latestPrice.Price, err
}

func fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string, onlyRecent bool) ([]models.BalanceResponse, error) {
	query := `
        SELECT
            b.id, b.ecdsa, b.eddsa, b.chain, b.address, b.token,
            b.balance, b.date, b.created_at, b.updated_at,
            COALESCE(b.balance * p.price, 0) as usd_value,
            p.id as price_id, p.chain as price_chain, p.token as price_token,
            p.price as price_price, p.date as price_date, p.source as price_source,
            p.created_at as price_created_at, p.updated_at as price_updated_at
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

	rows, err := db.DB.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []models.BalanceResponse
	for rows.Next() {
		var b models.BalanceResponse
		var p models.Price
		var priceID sql.NullInt64
		var usdValue sql.NullFloat64
		var priceChain, priceToken, priceSource sql.NullString
		var pricePrice sql.NullFloat64
		var priceDate sql.NullInt64
		var priceCreatedAt, priceUpdatedAt sql.NullTime

		err := rows.Scan(
			&b.ID, &b.ECDSA, &b.EDDSA, &b.Chain, &b.Address, &b.Token,
			&b.Balance, &b.Date, &b.CreatedAt, &b.UpdatedAt,
			&usdValue,
			&priceID, &priceChain, &priceToken, &pricePrice, &priceDate, &priceSource,
			&priceCreatedAt, &priceUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		b.UsdValue = usdValue.Float64

		if priceID.Valid {
			p.ID = uint(priceID.Int64)
			p.Chain = priceChain.String
			p.Token = priceToken.String
			p.Price = pricePrice.Float64
			p.Date = priceDate.Int64
			p.Source = priceSource.String
			p.CreatedAt = priceCreatedAt.Time
			p.UpdatedAt = priceUpdatedAt.Time
			b.Price = &p
		} else {
			recentPrice, err := GetLatestPriceByToken(b.Chain, b.Token)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			if recentPrice != nil {
				b.Price = recentPrice
				b.UsdValue = b.Balance * recentPrice.Price
			}
		}

		balances = append(balances, b)
	}

	return balances, nil
}

func GetBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.BalanceResponse, error) {
	return fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey, false)
}

func GetLatestBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.BalanceResponse, error) {
	return fetchBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey, true)
}

type ChainTokenPair struct {
	Chain string
	Token string
}

func GetUniqueActiveChainTokenPairs() ([]ChainTokenPair, error) {
	query := `
		SELECT DISTINCT chain, token
		FROM balances
		WHERE id IN (
			SELECT MAX(id)
			FROM balances
			GROUP BY chain, token, ecdsa, eddsa
		) AND balance > 0
	`

	rows, err := db.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pairs []ChainTokenPair
	for rows.Next() {
		var pair ChainTokenPair
		if err := rows.Scan(&pair.Chain, &pair.Token); err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func GetAverageBalanceSince(ecdsaPublicKey, eddsaPublicKey string, since time.Time) (float64, error) {
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

	var totalBalance float64
	err := db.DB.Raw(query, ecdsaPublicKey, eddsaPublicKey, since).Scan(&totalBalance).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get total average balance: %w", err)
	}

	return totalBalance, nil
}

func GetAverageBalanceForTimeRange(ecdsaPublicKey, eddsaPublicKey string, startTime, endTime time.Time) (float64, error) {
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

	var totalBalance float64
	err := db.DB.Raw(query, ecdsaPublicKey, eddsaPublicKey, startTime, endTime).Scan(&totalBalance).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get total average balance: %w", err)
	}

	return totalBalance, nil
}
