package services

import (
	"time"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func CreateCycle(cycle *models.Cycle) (*models.Cycle, error) {
	err := db.DB.Create(cycle).Error
	return cycle, err
}

func GetCurrentCycle() (*models.Cycle, error) {
	var cycle models.Cycle
	err := db.DB.Where("end_date > ?", time.Now()).Order("created_at desc").First(&cycle).Error
	return &cycle, err
}

func GetCycleByID(id uint) (*models.Cycle, error) {
	var cycle models.Cycle
	err := db.DB.Where("id = ?", id).First(&cycle).Error
	return &cycle, err
}

func GetAllCycles() ([]models.Cycle, error) {
	var cycles []models.Cycle
	err := db.DB.Order("created_at desc").Find(&cycles).Error
	return cycles, err
}
