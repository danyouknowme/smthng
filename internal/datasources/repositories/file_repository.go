package repositories

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/danyouknowme/smthng/internal/datasources"
)

type fileRepository struct {
	cld *cloudinary.Cloudinary
}

type FileRepository interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (string, error)
}

func NewFileRepository(ds datasources.DataSources) FileRepository {
	return &fileRepository{
		cld: ds.GetCloudinaryClient(),
	}
}

func (repo *fileRepository) Upload(ctx context.Context, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	uploadParam, err := repo.cld.Upload.Upload(ctx, f, uploader.UploadParams{
		Folder: "smthng",
	})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}
