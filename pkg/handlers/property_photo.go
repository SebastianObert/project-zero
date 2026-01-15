package handlers

import (
	"net/http"
	"project-zero/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PropertyPhotoHandler struct {
	db *gorm.DB
}

func NewPropertyPhotoHandler(db *gorm.DB) *PropertyPhotoHandler {
	return &PropertyPhotoHandler{db: db}
}

// AddPropertyPhoto menambahkan foto tambahan ke properti
func (h *PropertyPhotoHandler) AddPropertyPhoto(c *gin.Context) {
	var input models.PropertyPhoto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if err := h.db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// GetPropertyPhotos mengambil semua foto dari properti tertentu
func (h *PropertyPhotoHandler) GetPropertyPhotos(c *gin.Context) {
	propertyID := c.Param("property_id")

	var photos []models.PropertyPhoto
	if err := h.db.Where("property_id = ?", propertyID).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": photos})
}

// DeletePropertyPhoto menghapus foto tambahan
func (h *PropertyPhotoHandler) DeletePropertyPhoto(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.db.Delete(&models.PropertyPhoto{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
