package drivers

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/danyouknowme/smthng/pkg/logger"
)

func NewCloudinaryClient(cloudinaryURI string) (*cloudinary.Cloudinary, error) {
	logger.Info("Registering cloudinary client...")

	client, err := cloudinary.NewFromURL(cloudinaryURI)
	if err != nil {
		logger.Fatalf("Failed to create cloudinary client: %v", err)
		return nil, err
	}

	_, err = client.Admin.Ping(context.Background())
	if err != nil {
		logger.Fatalf("Failed to ping cloudinary: %v", err)
		return nil, err
	}

	logger.Info("Registering cloudinary client completed")
	return client, nil
}
