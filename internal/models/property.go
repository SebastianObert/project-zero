package models

import "time"

type Property struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ListingType string `json:"listing_type"` // WTS, WTR
	Price       int64  `json:"price"`

	// Detail Teknis
	LandSize     int `json:"land_size"`     // Luas Tanah (m2)
	BuildingSize int `json:"building_size"` // Luas Bangunan (m2)
	Bedrooms     int `json:"bedrooms"`
	Bathrooms    int `json:"bathrooms"`
	Floors       int `json:"floors"` // Jumlah Lantai

	// Fasilitas & Legalitas
	Certificate string `json:"certificate"`  // SHM, HGB, dll
	Electricity int    `json:"electricity"`  // Daya Listrik (Watt)
	WaterSource string `json:"water_source"` // PAM, Sumur
	Address     string `json:"address"`

	// Media
	PhotoPath string `json:"photo_path"` // Path foto properti

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
