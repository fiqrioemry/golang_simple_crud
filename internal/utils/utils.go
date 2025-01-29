package utils

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadMediaToCloudinary mengunggah file ke Cloudinary dengan kompresi
func UploadMediaToCloudinary(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	uploadFolder := os.Getenv("CLOUDINARY_FOLDER")

	// Inisialisasi Cloudinary
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Println("❌ Cloudinary configuration error:", err)
		return "", errors.New("cloudinary configuration failed")
	}

	// Upload file ke Cloudinary dengan transformasi untuk kompresi
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:         uploadFolder,
		Transformation: "w_500,h_500,c_fill,q_auto:low,f_auto", // Transformasi gambar
	})

	// Jika terjadi error dalam upload
	if err != nil {
		log.Println("❌ Failed to upload file to Cloudinary:", err)
		return "", errors.New("failed to upload file to Cloudinary")
	}

	// Berhasil, kembalikan URL file
	return uploadResult.SecureURL, nil
}
