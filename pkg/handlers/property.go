package handlers

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"project-zero/internal/models"
	"project-zero/pkg/database"
	"project-zero/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Constants untuk image validation
const (
	MaxFileSize   = 100 * 1024 * 1024 // 100MB
	MaxFileSizeMB = 100
)

var AllowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

type PropertyHandler struct {
	repo *database.PropertyRepository
}

// NewPropertyHandler membuat instance baru PropertyHandler
func NewPropertyHandler(repo *database.PropertyRepository) *PropertyHandler {
	return &PropertyHandler{repo: repo}
}

// CreateProperty membuat property baru
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	// Get userID from JWT middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input models.Property
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validasi gagal",
			"details": err.Error(),
		})
		return
	}

	// Set UserID dari token JWT
	input.UserID = userID.(uint)

	err := h.repo.CreateProperty(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal menyimpan data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// GetAllProperties mengambil semua property dengan pagination dan filtering
func (h *PropertyHandler) GetAllProperties(c *gin.Context) {
	// Get userID from JWT middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse query parameters
	params := utils.ParseQueryParams(c)

	// Tambahkan filter userID
	params.UserID = userID.(uint)

	// Fetch properties dari database (hanya milik user ini)
	properties, total, err := h.repo.GetPropertiesWithFilters(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal mengambil data",
			"details": err.Error(),
		})
		return
	}

	// Build pagination metadata
	pagination := utils.PaginationMetadata{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: utils.CalculateTotalPages(total, params.Limit),
	}

	// Build filters map untuk response
	filters := utils.BuildFiltersMap(params)

	// Return paginated response
	response := utils.PaginatedResponse{
		Data:       properties,
		Pagination: pagination,
	}
	if len(filters) > 0 {
		response.Filters = filters
	}

	c.JSON(http.StatusOK, response)
}

// GetPropertyByID mengambil property berdasarkan ID
func (h *PropertyHandler) GetPropertyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	property, err := h.repo.GetPropertyByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data gak ketemu"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Gagal mengambil data",
				"details": err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": property})
}

// UpdateProperty mengupdate property yang sudah ada
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input models.Property
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validasi gagal",
			"details": err.Error(),
		})
		return
	}

	property, err := h.repo.UpdateProperty(uint(id), &input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data gak ketemu"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Gagal mengupdate data",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": property})
}

// DeleteProperty menghapus property berdasarkan ID
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.repo.DeleteProperty(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal menghapus data",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Listing berhasil dihapus"})
}

// UploadFile menghandle upload file dengan validasi dan upload ke Cloudinary
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal ambil file"})
		return
	}

	// Validate file
	if err := validateImageFile(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Upload ke Cloudinary
	cloudinaryURL, err := utils.UploadToCloudinary(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal upload ke cloud",
			"details": err.Error(),
		})
		return
	}

	// Return Cloudinary URL untuk database
	c.JSON(http.StatusOK, gin.H{
		"photo_path": cloudinaryURL,
		"message":    "File berhasil diupload ke Cloudinary",
	})
}

// validateImageFile memvalidasi file image berdasarkan tipe dan ukuran
func validateImageFile(file *multipart.FileHeader) error {
	// Validasi ukuran file
	if file.Size > MaxFileSize {
		return &ValidationError{
			Code:    "FILE_TOO_LARGE",
			Message: "Ukuran file terlalu besar, maksimal " + strconv.Itoa(MaxFileSizeMB) + "MB",
		}
	}

	// Validasi tipe file (extension)
	ext := filepath.Ext(file.Filename)
	if !AllowedImageTypes[ext] {
		return &ValidationError{
			Code:    "INVALID_FILE_TYPE",
			Message: "Tipe file tidak didukung, gunakan: jpg, jpeg, png, gif, webp",
		}
	}

	return nil
}

// ValidationError custom error untuk validasi
type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
