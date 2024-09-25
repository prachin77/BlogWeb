package db

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	if err := godotenv.Load("P:/BlogWeb/.env"); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}

	Cloudinary_Cloud_Name := os.Getenv("Cloud_Name")
	Cloudinary_API_Key := os.Getenv("API_Key")
	Cloudinary_API_Secret := os.Getenv("API_Secret")

	cld, err := cloudinary.NewFromParams(Cloudinary_Cloud_Name, Cloudinary_API_Key, Cloudinary_API_Secret)
	if err != nil {
		return nil, fmt.Errorf("cloudinary setup error: %w", err)
	}
	return cld, nil
}

func UploadToCloudinary(file multipart.File, filename string) (string, string, error) {
	ctx := context.Background()
	cld, err := SetupCloudinary()
	if err != nil {
		return "", "", fmt.Errorf("error in creating upload instance: %w", err)
	}

	publicID := fmt.Sprintf("%d_%s", time.Now().Unix(), filename)

	uploadParams := uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: "image",
		Format:       "auto",
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", "", fmt.Errorf("error uploading to Cloudinary: %w", err)
	}

	return publicID, result.SecureURL, nil
}
