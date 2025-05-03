package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// UploadDir директория для хранения загруженных файлов
	UploadDir = "uploads"
	// MaxUploadSize максимальный размер загружаемого файла (10MB)
	MaxUploadSize = 10 << 20
)

// UploadImage загружает изображение и возвращает путь к файлу
func UploadImage(file *multipart.FileHeader) (string, error) {
	// Проверяем размер файла
	if file.Size > MaxUploadSize {
		return "", fmt.Errorf("размер файла превышает максимально допустимый (%d bytes)", MaxUploadSize)
	}

	// Проверяем тип файла (только изображения)
	if !isImageFile(file.Filename) {
		return "", fmt.Errorf("недопустимый формат файла (разрешены только jpg, jpeg, png, gif)")
	}

	// Создаем директорию для загрузок, если ее нет
	err := os.MkdirAll(UploadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("не удалось создать директорию для загрузок: %w", err)
	}

	// Создаем уникальное имя файла с использованием UUID и оригинального расширения
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102"), uuid.New().String(), ext)

	// Полный путь к новому файлу
	dst := filepath.Join(UploadDir, newFilename)

	// Открываем исходный файл
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть загружаемый файл: %w", err)
	}
	defer src.Close()

	// Создаем новый файл
	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("не удалось создать новый файл: %w", err)
	}
	defer out.Close()

	// Копируем содержимое
	_, err = io.Copy(out, src)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить файл: %w", err)
	}

	return dst, nil
}

// DeleteImage удаляет изображение по указанному пути
func DeleteImage(imagePath string) error {
	if imagePath == "" {
		return nil
	}

	if err := os.Remove(imagePath); err != nil {
		return fmt.Errorf("не удалось удалить файл: %w", err)
	}

	return nil
}

// isImageFile проверяет, является ли файл изображением по расширению
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}

	return false
}
