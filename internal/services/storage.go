package services

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/models"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	if nil == cfg {
		return nil, fmt.Errorf("config is nil")
	}
	mysqlConfig := cfg.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)

	database, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = database.AutoMigrate(&models.Vault{}, &models.CoinDBModel{}, &models.Job{}, &models.VaultShareAppearance{}, &models.VaultSeasonStats{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("connected to mysql database")
	return &Storage{db: database}, nil
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	return sqlDB.Close()
}

func (s *Storage) CreateJob(job *models.Job) error {
	result := s.db.Create(job)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetLastJob() (*models.Job, error) {
	var job models.Job
	result := s.db.Model(&models.Job{}).Last(&job)
	if result.Error != nil {
		return nil, result.Error
	}
	return &job, nil
}
func (s *Storage) UpdateJob(job *models.Job) error {
	result := s.db.Save(job)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateVaultRanks recalculates and updates the rank for all vaults with join_airdrop = 1,
// ensuring ranks are consecutive and sorted by total_points in descending order.
func (s *Storage) UpdateVaultRanks() error {
	sql := `
UPDATE vaults
    JOIN (
        SELECT id, ROW_NUMBER() OVER (ORDER BY total_points DESC) as vaultrank
        FROM vaults WHERE vaults.join_airdrop = 1
    ) ranked_vaults ON vaults.id = ranked_vaults.id
SET vaults.rank = ranked_vaults.vaultrank ;
`
	return s.db.Exec(sql).Error
}

func (s *Storage) UpdateVaultBalance() error {
	sql := `UPDATE vaults
		JOIN (
			SELECT vault_id, SUM(usd_value) AS total_balance
			FROM coins
			GROUP BY vault_id
		) AS coin_sums ON vaults.id = coin_sums.vault_id
		SET vaults.balance = coin_sums.total_balance;`
	return s.db.Exec(sql).Error
}

func (s *Storage) UpdateVaultTotalPoints() error {
	sql := `
        UPDATE vaults
        SET 
            total_points = total_points + SQRT(total_vault_value),
            total_vault_value = 0
    `
	return s.db.Exec(sql).Error
}
