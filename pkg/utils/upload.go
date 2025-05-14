package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UploadFile handles file uploads and returns the file path
func UploadFile(file *multipart.FileHeader, folderName string) (string, error) {
	// Create the uploads folder if it doesn't exist
	uploadsDir := filepath.Join("uploads", folderName)
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Error creating directory %s: %v", uploadsDir, err)
		return "", fmt.Errorf("error creating directory: %w", err)
	}

	// Generate a unique filename
	filename := generateUniqueFilename(file.Filename)

	// Full path for the file
	filepath := filepath.Join(uploadsDir, filename)

	// Create the file
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filepath, err)
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer dst.Close()

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Printf("Error opening uploaded file: %v", err)
		return "", fmt.Errorf("error opening uploaded file: %w", err)
	}
	defer src.Close()

	// Copy the file contents
	if _, err = io.Copy(dst, src); err != nil {
		log.Printf("Error copying file: %v", err)
		return "", fmt.Errorf("error copying file: %w", err)
	}

	// Return the path relative to the uploads directory
	return folderName + "/" + filename, nil
}

// generateUniqueFilename creates a unique filename for uploads
func generateUniqueFilename(originalFilename string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)

	// Get timestamp
	timestamp := time.Now().Format("20060102-150405")

	// Get UUID
	id := uuid.New().String()

	// Clean base name (without extension)
	baseName := strings.TrimSuffix(originalFilename, ext)
	baseName = strings.Replace(baseName, " ", "-", -1)
	baseName = strings.ToLower(baseName)

	// Return unique filename
	return fmt.Sprintf("%s-%s-%s%s", baseName, timestamp, id[:8], ext)
}

// DeleteFile removes a file from the uploads directory
func DeleteFile(filePath string) error {
	fullPath := filepath.Join("uploads", filePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// File does not exist, consider deletion successful
		log.Printf("File %s doesn't exist, skipping deletion", fullPath)
		return nil
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		log.Printf("Error deleting file %s: %v", fullPath, err)
		return fmt.Errorf("error deleting file: %w", err)
	}

	log.Printf("File %s successfully deleted", fullPath)
	return nil
}
