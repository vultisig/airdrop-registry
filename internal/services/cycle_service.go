package services

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func CreateCycle(cycle *models.Cycle) (*models.Cycle, error) {
	if err := db.DB.Create(cycle).Error; err != nil {
		return nil, fmt.Errorf("failed to create cycle: %w", err)
	}
	return cycle, nil
}

func GetCurrentCycle() (*models.Cycle, error) {
	var cycle models.Cycle
	if err := db.DB.Order("id desc").First(&cycle).Error; err != nil {
		return nil, fmt.Errorf("failed to get current cycle: %w", err)
	}
	return &cycle, nil
}

func GetCycleByID(id uint) (*models.Cycle, error) {
	var cycle models.Cycle
	if err := db.DB.First(&cycle, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get cycle with id %d: %w", id, err)
	}
	return &cycle, nil
}

func UpdateCycle(cycle *models.Cycle) error {
	if err := db.DB.Save(cycle).Error; err != nil {
		return fmt.Errorf("failed to update cycle: %w", err)
	}
	return nil
}

func DeleteCycle(id uint) error {
	if err := db.DB.Delete(&models.Cycle{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete cycle with id %d: %w", id, err)
	}
	return nil
}

func GetAllCycles(page, pageSize int) ([]models.Cycle, error) {
	var cycles []models.Cycle
	query := db.DB.Order("created_at DESC")

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&cycles).Error; err != nil {
		return nil, fmt.Errorf("failed to get all cycles: %w", err)
	}
	return cycles, nil
}
