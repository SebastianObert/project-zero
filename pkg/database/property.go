package database

import (
	"project-zero/internal/models"

	"gorm.io/gorm"
)

type PropertyRepository struct {
	db *gorm.DB
}

// NewPropertyRepository membuat instance baru PropertyRepository
func NewPropertyRepository(db *gorm.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

// CreateProperty menyimpan property baru ke database
func (r *PropertyRepository) CreateProperty(property *models.Property) error {
	result := r.db.Create(property)
	return result.Error
}

// GetAllProperties mengambil semua property dari database
func (r *PropertyRepository) GetAllProperties() ([]models.Property, error) {
	var properties []models.Property
	result := r.db.Find(&properties)
	return properties, result.Error
}

// GetPropertyByID mengambil property berdasarkan ID
func (r *PropertyRepository) GetPropertyByID(id uint) (*models.Property, error) {
	var property models.Property
	result := r.db.First(&property, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &property, nil
}

// UpdateProperty mengupdate property yang sudah ada
func (r *PropertyRepository) UpdateProperty(id uint, updates *models.Property) (*models.Property, error) {
	property, err := r.GetPropertyByID(id)
	if err != nil {
		return nil, err
	}

	result := r.db.Model(property).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	// Refresh data dari database
	r.db.First(property, id)
	return property, nil
}

// DeleteProperty menghapus property berdasarkan ID
func (r *PropertyRepository) DeleteProperty(id uint) error {
	result := r.db.Delete(&models.Property{}, id)
	return result.Error
}
