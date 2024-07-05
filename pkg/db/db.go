package db

import (
	"fmt"
	"log"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	mysqlConfig := config.Cfg.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = database.AutoMigrate(&models.Vault{}, &models.Balance{}, &models.Price{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("connected to mysql database")

	DB = database
}

func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from gorm.DB: %v", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("failed to close database: %v", err)
	}

	log.Println("disconnected from mysql database")
}
