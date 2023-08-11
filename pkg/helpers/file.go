package helpers

import (
	"mime/multipart"
	"net/http"
)

func GetFileContentType(out multipart.File) string {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return ""
	}

	contentType := http.DetectContentType(buffer)

	return contentType
}
