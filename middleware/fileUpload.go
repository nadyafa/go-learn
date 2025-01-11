package middleware

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

// max file size 10mb
const MaxFileSize = 10 * 1024 * 1024

var allowedFileTypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"application/pdf": ".pdf",
}

func CheckFileSize(file multipart.File) error {
	// read file into a buffer & determine its size
	buffer := make([]byte, MaxFileSize)
	_, err := file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("failed to read file")
	}

	// check file size limit
	if len(buffer) > MaxFileSize {
		return fmt.Errorf("file size must be less than 10 mb")
	}

	return nil
}

func CheckMimeType(ctx *gin.Context, file multipart.File) error {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file")
	}

	// detect MIME type
	mimeType := mimetype.Detect(buffer)

	// check if MIME type allowed
	if ext, ok := allowedFileTypes[mimeType.String()]; ok {
		ctx.Header("File-Extension", ext)
		return nil
	}

	// if MIME type not allowed
	return fmt.Errorf("invalid file type: %s", mimeType)
}

func GenerateFileName(projectName, originalFileName string) string {
	// get file extension from ori filename
	ext := filepath.Ext(originalFileName)

	// generate unique filename with username & timestamp 2006-01-02 15:04:05
	timestamp := time.Now().Format("20060102150405")
	newFileName := fmt.Sprintf("%s-%s%s", projectName, timestamp, ext)

	return newFileName
}
