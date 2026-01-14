package main

import (
	"fmt"
	"os"

	"project-zero/internal/models"
	"project-zero/pkg/database"
	"project-zero/pkg/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var propertyHandler *handlers.PropertyHandler

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

	// Initialize repository dan handler
	propertyRepo := database.NewPropertyRepository(db)
	propertyHandler = handlers.NewPropertyHandler(propertyRepo)

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
	r.POST("/upload", handlers.UploadFile)

	// Property routes
	r.POST("/properties", propertyHandler.CreateProperty)
	r.GET("/properties", propertyHandler.GetAllProperties)
	r.GET("/properties/:id", propertyHandler.GetPropertyByID)
	r.PUT("/properties/:id", propertyHandler.UpdateProperty)
	r.DELETE("/properties/:id", propertyHandler.DeleteProperty)

	r.Run(":8080")
}
