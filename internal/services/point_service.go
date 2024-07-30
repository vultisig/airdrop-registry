package services

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
	"gorm.io/gorm"
)

func SavePoint(point *models.Point) error {
	var existingPoint models.Point
	err := db.DB.Where("ecdsa = ? AND eddsa = ? AND cycle_id = ?", point.ECDSA, point.EDDSA, point.CycleID).First(&existingPoint).Error
	if err == nil {
		existingPoint.Balance = point.Balance
		existingPoint.Share = point.Share
		if err := db.DB.Save(&existingPoint).Error; err != nil {
			return fmt.Errorf("failed to update existing point: %w", err)
		}
		return nil
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing point: %w", err)
	}

	if err := db.DB.Create(point).Error; err != nil {
		return fmt.Errorf("failed to create new point: %w", err)
	}
	return nil
}

func GetPointByID(id uint) (*models.Point, error) {
	var point models.Point
	if err := db.DB.Joins("Cycle").First(&point, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get point with id %d: %w", id, err)
	}
	return &point, nil
}

func UpdatePoint(point *models.Point) error {
	if err := db.DB.Save(point).Error; err != nil {
		return fmt.Errorf("failed to update point: %w", err)
	}
	return nil
}

func DeletePoint(id uint) error {
	if err := db.DB.Delete(&models.Point{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete point with id %d: %w", id, err)
	}
	return nil
}

func GetAllPoints(page, pageSize int) ([]models.Point, error) {
	var points []models.Point
	query := db.DB.Joins("Cycle").Order("points.created_at DESC")

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&points).Error; err != nil {
		return nil, fmt.Errorf("failed to get all points: %w", err)
	}
	return points, nil
}

func GetLatestPointByVault(ecdsaPublicKey, eddsaPublicKey string) (*models.Point, error) {
	var point models.Point
	err := db.DB.Joins("Cycle").
		Where("points.ecdsa = ? AND points.eddsa = ?", ecdsaPublicKey, eddsaPublicKey).
		Order("points.created_at desc").
		First(&point).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get latest point for vault: %w", err)
	}
	return &point, nil
}

func GetPointsByVault(ecdsaPublicKey, eddsaPublicKey string) ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Joins("Cycle").
		Where("points.ecdsa = ? AND points.eddsa = ?", ecdsaPublicKey, eddsaPublicKey).
		Order("points.created_at desc").
		Find(&points).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get points for vault: %w", err)
	}
	return points, nil
}

func GetPointsForVaultByCycle(ecdsaPublicKey, eddsaPublicKey string, cycleID uint) ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Joins("Cycle").
		Where("points.ecdsa = ? AND points.eddsa = ? AND points.cycle_id = ?", ecdsaPublicKey, eddsaPublicKey, cycleID).
		Find(&points).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get points for vault by cycle: %w", err)
	}
	return points, nil
}
