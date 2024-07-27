package services

import (
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
		return db.DB.Save(&existingPoint).Error
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return db.DB.Create(point).Error
}

func GetLatestPointByVault(ecdsaPublicKey, eddsaPublicKey string) (*models.Point, error) {
	var point models.Point
	err := db.DB.Joins("Cycle").
		Where("points.ecdsa = ? AND points.eddsa = ?", ecdsaPublicKey, eddsaPublicKey).
		Order("points.created_at desc").
		First(&point).Error
	return &point, err
}

func GetPointsByVault(ecdsaPublicKey, eddsaPublicKey string) ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Joins("Cycle").
		Where("points.ecdsa = ? AND points.eddsa = ?", ecdsaPublicKey, eddsaPublicKey).
		Order("points.created_at desc").
		Find(&points).Error
	return points, err
}

func GetAllPoints() ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Joins("Cycle").
		Order("points.created_at desc").
		Find(&points).Error
	return points, err
}

func GetPointByID(id uint) (*models.Point, error) {
	var point models.Point
	err := db.DB.Joins("Cycle").
		Where("points.id = ?", id).
		First(&point).Error
	return &point, err
}

func GetPointsForVaultByCycle(ecdsaPublicKey, eddsaPublicKey string, cycleID uint) ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Joins("Cycle").
		Where("points.ecdsa = ? AND points.eddsa = ? AND points.cycle_id = ?", ecdsaPublicKey, eddsaPublicKey, cycleID).
		Find(&points).Error
	return points, err
}
