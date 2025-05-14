package service

import (
	"FurniSwap/internal/modules/profile/model"
	"FurniSwap/internal/modules/profile/repository"
	"FurniSwap/pkg/utils"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Service provides profile operations
type Service struct {
	repo *repository.Repository
}

// NewService creates a new profile service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetProfile retrieves a user profile
func (s *Service) GetProfile(userID int) (*model.Profile, error) {
	return s.repo.GetProfileByID(userID)
}

// UpdateProfile updates a user profile
func (s *Service) UpdateProfile(userID int, req model.UpdateProfileRequest) error {
	// Если задан URL аватара, проверяем его и обрабатываем
	if req.Avatar != "" && (strings.HasPrefix(req.Avatar, "http://") || strings.HasPrefix(req.Avatar, "https://")) {
		log.Printf("Updating profile with avatar URL: %s", req.Avatar)

		// Проверяем URL аватара по тем же правилам, что и в SetAvatarURL
		// Базовая проверка на URL уже пройдена
		urlLower := strings.ToLower(req.Avatar)
		isImageByExtension := strings.HasSuffix(urlLower, ".jpg") ||
			strings.HasSuffix(urlLower, ".jpeg") ||
			strings.HasSuffix(urlLower, ".png") ||
			strings.HasSuffix(urlLower, ".gif") ||
			strings.HasSuffix(urlLower, ".webp") ||
			strings.HasSuffix(urlLower, ".bmp") ||
			strings.HasSuffix(urlLower, ".svg")

		// Если это не явно изображение по расширению, проверяем URL мягко (без возврата ошибки)
		if !isImageByExtension {
			client := &http.Client{
				Timeout: 10 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					if len(via) >= 10 {
						return fmt.Errorf("слишком много редиректов")
					}
					return nil
				},
			}

			req, err := http.NewRequest("HEAD", req.Avatar, nil)
			if err == nil {
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

				resp, err := client.Do(req)
				if err == nil {
					defer resp.Body.Close()

					contentType := resp.Header.Get("Content-Type")
					if contentType != "" && !strings.Contains(contentType, "image/") {
						log.Printf("Warning: Avatar URL не имеет Content-Type изображения: %s", contentType)
					}
				}
			}
		}

		// Удаляем старый локальный аватар, если он был
		currentAvatar, err := s.repo.GetAvatarPath(userID)
		if err == nil && currentAvatar != "" && !s.IsAvatarURL(currentAvatar) {
			_ = utils.DeleteFile(currentAvatar)
		}
	}

	// Обновляем профиль
	return s.repo.UpdateProfile(userID, req)
}

// UploadAvatar uploads a user avatar
func (s *Service) UploadAvatar(userID int, file *multipart.FileHeader) error {
	// Check if user has an existing avatar
	currentAvatar, err := s.repo.GetAvatarPath(userID)
	if err != nil {
		log.Printf("Error getting current avatar: %v", err)
	}

	// If there's an existing avatar and it's not a URL, delete it
	if currentAvatar != "" && !s.IsAvatarURL(currentAvatar) {
		err = utils.DeleteFile(currentAvatar)
		if err != nil {
			log.Printf("Error deleting old avatar: %v", err)
		}
	}

	// Upload the new avatar
	avatarPath, err := utils.UploadFile(file, "avatars")
	if err != nil {
		return fmt.Errorf("error uploading avatar: %w", err)
	}

	// Update the user's avatar in the database
	err = s.repo.UpdateAvatar(userID, avatarPath)
	if err != nil {
		// If there's an error updating the database, try to clean up the uploaded file
		_ = utils.DeleteFile(avatarPath)
		return fmt.Errorf("error updating avatar in database: %w", err)
	}

	return nil
}

// SetAvatarURL sets a URL as the user's avatar
func (s *Service) SetAvatarURL(userID int, url string) error {
	// Basic validation of the URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("invalid image URL, must start with http:// or https://")
	}

	// Проверка на наличие расширения файла в URL, характерного для изображений
	urlLower := strings.ToLower(url)
	isImageByExtension := strings.HasSuffix(urlLower, ".jpg") ||
		strings.HasSuffix(urlLower, ".jpeg") ||
		strings.HasSuffix(urlLower, ".png") ||
		strings.HasSuffix(urlLower, ".gif") ||
		strings.HasSuffix(urlLower, ".webp") ||
		strings.HasSuffix(urlLower, ".bmp") ||
		strings.HasSuffix(urlLower, ".svg")

	// Если URL явно указывает на изображение по расширению, пропускаем проверку HEAD запросом
	if !isImageByExtension {
		// Проверка URL с помощью HEAD-запроса
		client := &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Разрешаем до 10 редиректов
				if len(via) >= 10 {
					return fmt.Errorf("слишком много редиректов")
				}
				return nil
			},
		}

		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		// Добавление User-Agent для обхода блокировок
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error checking URL %s: %v", url, err)
			// Если не смогли проверить URL, все равно попробуем его использовать
			// потому что многие сайты блокируют HEAD-запросы
		} else {
			defer resp.Body.Close()

			// Проверяем статус-код
			if resp.StatusCode >= 400 {
				return fmt.Errorf("URL недоступен, статус: %d", resp.StatusCode)
			}

			// Проверяем Content-Type только если он есть
			contentType := resp.Header.Get("Content-Type")
			if contentType != "" && !strings.Contains(contentType, "image/") {
				log.Printf("Warning: URL не имеет Content-Type изображения: %s, но попробуем использовать", contentType)
				// Не возвращаем ошибку, т.к. контент может быть изображением, несмотря на неправильный Content-Type
			}
		}
	}

	// Check if user has an existing avatar
	currentAvatar, err := s.repo.GetAvatarPath(userID)
	if err != nil {
		log.Printf("Error getting current avatar: %v", err)
	}

	// If there's an existing avatar and it's not a URL, delete it
	if currentAvatar != "" && !s.IsAvatarURL(currentAvatar) {
		err = utils.DeleteFile(currentAvatar)
		if err != nil {
			log.Printf("Error deleting old avatar: %v", err)
		}
	}

	// Update the user's avatar in the database
	err = s.repo.UpdateAvatar(userID, url)
	if err != nil {
		return fmt.Errorf("error updating avatar URL in database: %w", err)
	}

	return nil
}

// IsAvatarURL checks if the avatar path is a URL
func (s *Service) IsAvatarURL(avatarPath string) bool {
	return strings.HasPrefix(avatarPath, "http://") || strings.HasPrefix(avatarPath, "https://")
}

// GetAvatarPath gets the avatar path for a user
func (s *Service) GetAvatarPath(userID int) (string, error) {
	return s.repo.GetAvatarPath(userID)
}

// GetPublicProfile converts a profile to a public profile
func (s *Service) GetPublicProfile(profile *model.Profile) *model.PublicProfile {
	return &model.PublicProfile{
		ID:        profile.ID,
		Name:      profile.Name,
		LastName:  profile.LastName,
		City:      profile.City,
		Avatar:    profile.Avatar,
		CreatedAt: profile.CreatedAt,
	}
}
