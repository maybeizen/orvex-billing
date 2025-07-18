package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxAvatarSize = 5 * 1024 * 1024 // 5MB
	AvatarWidth   = 400
	AvatarHeight  = 400
)

var AllowedAvatarTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/webp": true,
}

type FileValidationError struct {
	Message string
}

func (e FileValidationError) Error() string {
	return e.Message
}

func ValidateAvatar(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > MaxAvatarSize {
		return FileValidationError{"File size exceeds 5MB limit"}
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
		return FileValidationError{"Only PNG, JPG, and WebP files are allowed"}
	}

	file, err := fileHeader.Open()
	if err != nil {
		return FileValidationError{"Failed to open file"}
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return FileValidationError{"Failed to read file"}
	}

	file.Seek(0, 0)

	contentType := DetectContentType(buffer[:n])
	if !AllowedAvatarTypes[contentType] {
		return FileValidationError{"Invalid file type detected"}
	}

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return FileValidationError{"Invalid image format"}
	}

	switch format {
	case "png", "jpeg", "webp":
	default:
		return FileValidationError{"Unsupported image format"}
	}

	return nil
}

func DetectContentType(data []byte) string {
	if len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png"
	}

	// Improved JPEG detection - just check for SOI marker (0xFF 0xD8)
	// The third byte can be any marker (0xE0, 0xE1, 0xDB, etc.)
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	if len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}

	return "unknown"
}

func GenerateSecureFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes) + ext
}

func SaveAvatar(fileHeader *multipart.FileHeader, userID uint) (string, error) {
	if err := ValidateAvatar(fileHeader); err != nil {
		return "", err
	}

	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	filename := GenerateSecureFilename(fileHeader.Filename)
	filepath := fmt.Sprintf("%s/%d_%s", uploadDir, userID, filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filepath) // Clean up on error
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return fmt.Sprintf("/uploads/avatars/%d_%s", userID, filename), nil
}

func DeleteAvatar(avatarURL string) error {
	if avatarURL == "" {
		return nil
	}

	filepath := "." + avatarURL
	
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filepath)
}
