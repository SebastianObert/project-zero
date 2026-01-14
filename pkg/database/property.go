package database

import (
	"project-zero/internal/models"
	"project-zero/pkg/utils"

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

// GetPropertiesWithFilters mengambil properties dengan pagination, filtering, dan sorting
func (r *PropertyRepository) GetPropertiesWithFilters(params utils.QueryParams) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	query := r.db

	// Apply filters
	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}
	if params.ListingType != "" {
		query = query.Where("listing_type = ?", params.ListingType)
	}
	if params.Bedrooms > 0 {
		query = query.Where("bedrooms = ?", params.Bedrooms)
	}
	if params.Bathrooms > 0 {
		query = query.Where("bathrooms = ?", params.Bathrooms)
	}
	if params.Certificate != "" {
		query = query.Where("certificate = ?", params.Certificate)
	}
	if params.Location != "" {
		// Use ILIKE for case-insensitive search
		query = query.Where("address ILIKE ?", "%"+params.Location+"%")
	}
	if params.Title != "" {
		query = query.Where("title ILIKE ?", "%"+params.Title+"%")
	}

	// Get total count BEFORE pagination
	countResult := query.Model(&models.Property{}).Count(&total)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	// Apply sorting
	sortColumn := params.SortBy
	sortOrder := params.SortOrder
	query = query.Order(sortColumn + " " + sortOrder)

	// Apply pagination
	offset := utils.CalculateOffset(params.Page, params.Limit)
	result := query.Offset(offset).Limit(params.Limit).Find(&properties)

	return properties, total, result.Error
}
