package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds the application configuration
type AppConfig struct {
	// Server settings
	Port string

	// CORS settings
	AllowedOrigins []string

	// JWT Settings
	JWTSecret string

	// Database settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// SMTP settings
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string

	// File upload settings
	UploadsDir string
}

// Config is the global application configuration
var Config AppConfig

// Load loads the application configuration from environment variables
func Load() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file:", err)
		// Try loading env.txt
		err = godotenv.Load("env.txt")
		if err != nil {
			log.Println("Warning: Error loading env.txt file. Using default values.")
		}
	}

	// Server settings
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// CORS settings
	allowedOrigins := []string{"http://localhost:3000"}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = []string{origins}
	}

	// JWT Settings
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key_change_in_production"
		log.Println("WARNING: JWT_SECRET_KEY not configured, using default value")
	}

	// Database settings
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5431"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "furni_swap"
	}

	// SMTP settings
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// File upload settings
	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "uploads"
	}

	// Create the uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatal("Error creating uploads directory:", err)
	}

	// Set the global configuration
	Config = AppConfig{
		Port:           port,
		AllowedOrigins: allowedOrigins,
		JWTSecret:      jwtSecret,
		DBHost:         dbHost,
		DBPort:         dbPort,
		DBUser:         dbUser,
		DBPassword:     dbPassword,
		DBName:         dbName,
		SMTPHost:       smtpHost,
		SMTPPort:       smtpPort,
		SMTPUsername:   smtpUsername,
		SMTPPassword:   smtpPassword,
		UploadsDir:     uploadsDir,
	}

	log.Println("Configuration loaded successfully")
}

// GetConfig returns the application configuration
func GetConfig() AppConfig {
	return Config
}
