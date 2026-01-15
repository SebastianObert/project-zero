package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

// InitCloudinary inisialisasi Cloudinary client
func InitCloudinary() error {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return fmt.Errorf("Cloudinary credentials belum diatur di .env")
	}

	var err error
	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("Gagal inisialisasi Cloudinary: %v", err)
	}

	fmt.Println("✅ Cloudinary berhasil diinisialisasi!")
	return nil
}

// UploadToCloudinary upload file ke Cloudinary dan return URL-nya
func UploadToCloudinary(file *multipart.FileHeader) (string, error) {
	if cld == nil {
		return "", fmt.Errorf("Cloudinary belum diinisialisasi")
	}

	// Buka file untuk dibaca
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("Gagal membuka file: %v", err)
	}
	defer src.Close()

	// Generate unique public ID dengan timestamp
	ext := filepath.Ext(file.Filename)
	publicID := fmt.Sprintf("properties/%s_%s", time.Now().Format("20060102150405"), file.Filename[:len(file.Filename)-len(ext)])

	// Upload ke Cloudinary
	ctx := context.Background()
	overwrite := false
	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: "image",
		Folder:       "property-photos", // Semua foto property masuk ke folder ini
		Overwrite:    &overwrite,
	})

	if err != nil {
		return "", fmt.Errorf("Gagal upload ke Cloudinary: %v", err)
	}

	// Return secure URL (HTTPS)
	return uploadResult.SecureURL, nil
}

// DeleteFromCloudinary menghapus file dari Cloudinary berdasarkan public ID
func DeleteFromCloudinary(publicID string) error {
	if cld == nil {
		return fmt.Errorf("Cloudinary belum diinisialisasi")
	}

	ctx := context.Background()
	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return fmt.Errorf("Gagal hapus dari Cloudinary: %v", err)
	}

	return nil
}
