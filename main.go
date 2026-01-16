package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"project-zero/internal/models"
	"project-zero/pkg/database"
	"project-zero/pkg/handlers"
	"project-zero/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var propertyHandler *handlers.PropertyHandler
var propertyPhotoHandler *handlers.PropertyPhotoHandler
var authHandler *handlers.AuthHandler

func initDB() {
	// Ambil config dari .env
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSL_MODE"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("❌ Gagal koneksi ke database: %v", err))
	}
	fmt.Println("✅ Berhasil terhubung ke Neon Database!")

	// Buat/update tabel otomatis
	db.AutoMigrate(&models.User{}, &models.Property{}, &models.PropertyPhoto{})

	// Initialize Cloudinary
	if err := utils.InitCloudinary(); err != nil {
		fmt.Printf("⚠️  Warning: %v\n", err)
		fmt.Println("📝 Silakan isi CLOUDINARY credentials di file .env")
	}

	// Initialize repository dan handler
	propertyRepo := database.NewPropertyRepository(db)
	propertyHandler = handlers.NewPropertyHandler(propertyRepo)
	propertyPhotoHandler = handlers.NewPropertyPhotoHandler(db)
	authHandler = handlers.NewAuthHandler(db)

	// Buat folder untuk upload foto jika belum ada (optional, backup only)
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

	// Set max multipart memory untuk support file besar (50MB lebih aman)
	r.MaxMultipartMemory = 50 * 1024 * 1024 // 50MB

	// Pasang middleware CORS biar Frontend (HTML) bisa akses Backend
	r.Use(corsMiddleware())

	// Serve static files (HTML, CSS, JS)
	r.StaticFile("/", "./index.html")
	r.StaticFile("/index.html", "./index.html")
	r.StaticFile("/login.html", "./login.html")
	r.StaticFile("/signup.html", "./signup.html")
	r.StaticFile("/loading-modal.html", "./loading-modal.html")
	r.StaticFile("/styles.css", "./styles.css")
	r.StaticFile("/auth-helper.js", "./auth-helper.js")

	// Auth routes (PUBLIC - tidak perlu login)
	auth := r.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes (PRIVATE - perlu login dengan JWT)
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		// Profile
		protected.GET("/auth/profile", authHandler.GetProfile)

		// Upload foto endpoint - upload langsung ke Cloudinary
		protected.POST("/upload", handlers.UploadFile)

		// Property routes
		protected.POST("/properties", propertyHandler.CreateProperty)
		protected.GET("/properties", propertyHandler.GetAllProperties)
		protected.GET("/properties/:id", propertyHandler.GetPropertyByID)
		protected.PUT("/properties/:id", propertyHandler.UpdateProperty)
		protected.DELETE("/properties/:id", propertyHandler.DeleteProperty)

		// Property photos routes
		protected.POST("/property-photos", propertyPhotoHandler.AddPropertyPhoto)
		protected.GET("/property-photos/:property_id", propertyPhotoHandler.GetPropertyPhotos)
		protected.DELETE("/property-photos/:id", propertyPhotoHandler.DeletePropertyPhoto)
	}

	// Create HTTP server with custom timeouts
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    90 * time.Second, // Timeout untuk baca request
		WriteTimeout:   90 * time.Second, // Timeout untuk kirim response
		MaxHeaderBytes: 1 << 20,          // 1 MB
	}

	fmt.Println("🚀 Server berjalan di http://localhost:8080")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("❌ Gagal start server: %v", err))
	}
}
