package database

import (
	"project-zero/internal/models"
	"project-zero/pkg/utils"

	"gorm.io/gorm"
)

type PropertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

func (r *PropertyRepository) CreateProperty(property *models.Property) error {
	return r.db.Create(property).Error
}

func (r *PropertyRepository) GetPropertiesWithFilters(params utils.QueryParams) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	query := r.db.Model(&models.Property{})

	// Filter berdasarkan UserID
	if params.UserID > 0 {
		query = query.Where("user_id = ?", params.UserID)
	}

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
		query = query.Where("bedrooms >= ?", params.Bedrooms)
	}

	// Count total
	query.Count(&total)

	// Apply sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := params.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	err := query.Limit(params.Limit).Offset(offset).Find(&properties).Error

	return properties, total, err
}
func (r *PropertyRepository) GetPropertyByID(id uint) (*models.Property, error) {
	var property models.Property
	err := r.db.First(&property, id).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}

func (r *PropertyRepository) UpdateProperty(id uint, property *models.Property) (*models.Property, error) {
	var existing models.Property
	if err := r.db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	// Update all fields using Updates with map to handle zero values
	updates := map[string]interface{}{
		"title":         property.Title,
		"description":   property.Description,
		"price":         property.Price,
		"listing_type":  property.ListingType,
		"land_size":     property.LandSize,
		"building_size": property.BuildingSize,
		"bedrooms":      property.Bedrooms,
		"bathrooms":     property.Bathrooms,
		"floors":        property.Floors,
		"certificate":   property.Certificate,
		"electricity":   property.Electricity,
		"water_source":  property.WaterSource,
		"address":       property.Address,
		"photo_path":    property.PhotoPath,
	}

	if err := r.db.Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Fetch the updated record
	if err := r.db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (r *PropertyRepository) DeleteProperty(id uint) error {
	return r.db.Delete(&models.Property{}, id).Error
}
