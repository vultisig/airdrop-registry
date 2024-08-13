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
			Logger: logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = database.AutoMigrate(&models.Vault{}, &models.CoinDBModel{})
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
