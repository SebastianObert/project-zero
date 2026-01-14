package models

import "time"

type Property struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ListingType  string    `json:"listing_type"` // WTS, WTB, WTR
	Price        int64     `json:"price"`
	LandSize     int       `json:"land_size"`     // Luas Tanah
	BuildingSize int       `json:"building_size"` // Luas Bangunan
	Bedrooms     int       `json:"bedrooms"`
	Bathrooms    int       `json:"bathrooms"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
