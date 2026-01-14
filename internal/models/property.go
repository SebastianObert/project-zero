package models

import "time"

type Property struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required,min=10,max=2000"`
	ListingType string `json:"listing_type" binding:"required,oneof=WTS WTR"` // WTS, WTR
	Price       int64  `json:"price" binding:"required,gt=0"`

	// Detail Teknis
	LandSize     int `json:"land_size" binding:"required,gt=0"`     // Luas Tanah (m2)
	BuildingSize int `json:"building_size" binding:"required,gt=0"` // Luas Bangunan (m2)
	Bedrooms     int `json:"bedrooms" binding:"required,gte=0,lte=20"`
	Bathrooms    int `json:"bathrooms" binding:"required,gte=0,lte=20"`
	Floors       int `json:"floors" binding:"required,gte=1,lte=50"` // Jumlah Lantai

	// Fasilitas & Legalitas
	Certificate string `json:"certificate" binding:"required,oneof=SHM HGB GIRIK LAINNYA"`  // SHM, HGB, dll
	Electricity int    `json:"electricity" binding:"required,gt=0"`                         // Daya Listrik (Watt)
	WaterSource string `json:"water_source" binding:"required,oneof=PAM Sumur Bor LAINNYA"` // PAM, Sumur
	Address     string `json:"address" binding:"required,min=10,max=500"`

	// Media
	PhotoPath string `json:"photo_path"` // Path foto properti

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
