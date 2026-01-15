package models

import "time"

type PropertyPhoto struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PropertyID uint      `json:"property_id" binding:"required"`
	PhotoPath  string    `json:"photo_path" binding:"required"`
	Caption    string    `json:"caption"`
	CreatedAt  time.Time `json:"created_at"`
}
