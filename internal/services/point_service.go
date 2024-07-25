package services

import (
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func SavePoint(point *models.Point) error {
	return db.DB.Create(point).Error
}

func GetLatestPointByVault(ecdsaPublicKey, eddsaPublicKey string) (*models.Point, error) {
	var point models.Point
	err := db.DB.Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).Order("created_at desc").First(&point).Error
	return &point, err
}

func GetPointsByVault(ecdsaPublicKey, eddsaPublicKey string) ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).Order("created_at desc").Find(&points).Error
	return points, err
}

func GetAllPoints() ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Order("created_at desc").Find(&points).Error
	return points, err
}

func GetPointByID(id uint) (*models.Point, error) {
	var point models.Point
	err := db.DB.Where("id = ?", id).First(&point).Error
	return &point, err
}

func GetPointsForVaultByTokenAndDateRange(ecdsaPublicKey, eddsaPublicKey, chain, token string, start, end int64) ([]models.Point, error) {
	var points []models.Point
	err := db.DB.Where("ecdsa = ? AND eddsa = ? AND chain = ? AND token = ? AND date >= ? AND date <= ?", ecdsaPublicKey, eddsaPublicKey, chain, token, start, end).Find(&points).Error
	return points, err
}
