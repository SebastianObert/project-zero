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

// Middleware untuk HTTPS redirect di production
func httpsRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Hanya redirect di production
		if os.Getenv("ENVIRONMENT") == "production" {
			// Check X-Forwarded-Proto header (untuk proxy seperti nginx/cloudflare)
			proto := c.GetHeader("X-Forwarded-Proto")
			if proto == "http" {
				httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
				c.Redirect(http.StatusMovedPermanently, httpsURL)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// Middleware untuk security headers
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// HTTPS Strict Transport Security (HSTS)
		if os.Getenv("ENVIRONMENT") == "production" {
			c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Security headers
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (CSP)
		if os.Getenv("ENVIRONMENT") == "production" {
			c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https://res.cloudinary.com; connect-src 'self';")
		}

		c.Next()
	}
}

// Middleware buat handle CORS (Cross-Origin Resource Sharing)
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Di production, hanya allow domain tertentu
		if os.Getenv("ENVIRONMENT") == "production" {
			allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
			if allowedOrigins != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			// Development: allow all
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	initDB()

	// Set Gin mode berdasarkan environment
	if os.Getenv("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Set max multipart memory untuk support file besar
	r.MaxMultipartMemory = 50 * 1024 * 1024 // 50MB

	// Security middlewares (urutan penting!)
	r.Use(securityHeadersMiddleware())
	r.Use(httpsRedirectMiddleware())
	r.Use(corsMiddleware())

	// Serve static files (HTML, CSS, JS)
	r.StaticFile("/", "./index.html")
	r.StaticFile("/index.html", "./index.html")
	r.StaticFile("/login.html", "./login.html")
	r.StaticFile("/signup.html", "./signup.html")
	r.StaticFile("/loading-modal.html", "./loading-modal.html")
	r.StaticFile("/styles.css", "./styles.css")
	r.StaticFile("/auth-helper.js", "./auth-helper.js")

	// Health check endpoint (untuk monitoring & load balancer)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"environment": os.Getenv("ENVIRONMENT"),
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		})
	})

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

	// Get port dari environment atau default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server with custom timeouts
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    90 * time.Second, // Timeout untuk baca request
		WriteTimeout:   90 * time.Second, // Timeout untuk kirim response
		MaxHeaderBytes: 1 << 20,          // 1 MB
	}

	// Start server based on environment
	if os.Getenv("ENVIRONMENT") == "production" {
		// Production: HTTPS
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")

		if certFile != "" && keyFile != "" {
			fmt.Printf("🔒 Server berjalan di https://0.0.0.0:%s (PRODUCTION - HTTPS)\n", port)
			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("❌ Gagal start HTTPS server: %v", err))
			}
		} else {
			fmt.Printf("⚠️  Production mode tetapi TLS tidak dikonfigurasi\n")
			fmt.Printf("🚀 Server berjalan di http://0.0.0.0:%s (WARNING: HTTP in production)\n", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("❌ Gagal start server: %v", err))
			}
		}
	} else {
		// Development: HTTP
		fmt.Printf("🚀 Server berjalan di http://localhost:%s (DEVELOPMENT)\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("❌ Gagal start server: %v", err))
		}
	}
}
