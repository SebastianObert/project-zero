package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"project-zero/internal/models"
	"project-zero/pkg/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var propertyRepo *database.PropertyRepository

func initDB() {
	// Ambil config dari .env
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke database")
	}

	// Buat/update tabel otomatis
	db.AutoMigrate(&models.Property{})

	// Initialize repository
	propertyRepo = database.NewPropertyRepository(db)

	// Buat folder untuk upload foto jika belum ada
	os.MkdirAll("./uploads", os.ModePerm)
}

// Middleware buat handle CORS (Cross-Origin Resource Sharing)
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	initDB()

	r := gin.Default()

	// Pasang middleware CORS biar Frontend (HTML) bisa akses Backend
	r.Use(corsMiddleware())

	// Serve upload folder sebagai static files
	r.Static("/uploads", "./uploads")

	// Upload foto endpoint
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal ambil file"})
			return
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		uploadPath := "./uploads/" + filename

		// Save file
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan file"})
			return
		}

		// Return URL path untuk database
		c.JSON(http.StatusOK, gin.H{"photo_path": "/uploads/" + filename})
	})

	// 1. Create: Simpan listing baru
	r.POST("/properties", func(c *gin.Context) {
		var input models.Property
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validasi gagal",
				"details": err.Error(),
			})
			return
		}

		err := propertyRepo.CreateProperty(&input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Gagal menyimpan data",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": input})
	})

	// 2. Read All: Ambil semua list rumah
	r.GET("/properties", func(c *gin.Context) {
		properties, err := propertyRepo.GetAllProperties()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Gagal mengambil data",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": properties})
	})

	// 3. Read One: Cari rumah spesifik pake ID
	r.GET("/properties/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		property, err := propertyRepo.GetPropertyByID(uint(id))
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
	})

	// 4. Update: Edit data rumah
	r.PUT("/properties/:id", func(c *gin.Context) {
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

		property, err := propertyRepo.UpdateProperty(uint(id), &input)
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
	})

	// 5. Delete: Hapus listing
	r.DELETE("/properties/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		err = propertyRepo.DeleteProperty(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Gagal menghapus data",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Listing berhasil dihapus"})
	})

	r.Run(":8080")
}
