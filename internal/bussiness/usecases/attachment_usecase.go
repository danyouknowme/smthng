package usecases

import (
	"mime/multipart"

	"github.com/danyouknowme/smthng/internal/datasources/repositories"
)

type attachmentUsecase struct {
	repo repositories.FileRepository
}

type AttachmentUsecase interface {
	UploadFile(fileHeader *multipart.FileHeader) (string, error)
}

func NewAttachmentUsecase(repo repositories.FileRepository) AttachmentUsecase {
	return &attachmentUsecase{
		repo: repo,
	}
}

func (uc *attachmentUsecase) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	// uploadObjectID, err := uc.repo.Upload(fileHeader)
	// if err != nil {
	// 	return "", err
	// }

	return "", nil
}
