package database

import (
	"FurniSwap/pkg/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// DB is the global database instance
var DB *sqlx.DB

// Init initializes the database connection
func Init() (*sqlx.DB, error) {
	// Build connection string from config
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config.DBHost, config.Config.DBPort, config.Config.DBUser, config.Config.DBPassword, config.Config.DBName)
	log.Println("Connecting to database:", connectionString)

	// Connect to the database
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Set the global DB instance
	DB = db

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// GetDB returns the global database instance
func GetDB() *sqlx.DB {
	return DB
}

// Close closes the database connection
func Close() {
	if DB != nil {
		log.Println("Closing database connection")
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

// UpdateDatabaseStructure checks and updates database structure if needed
func UpdateDatabaseStructure(db *sqlx.DB) {
	// Check structure of users table
	var columns []string
	err := db.Select(&columns, `
		SELECT column_name FROM information_schema.columns 
		WHERE table_name = 'users'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Printf("Error checking table structure: %v", err)
		return
	}

	log.Printf("Current users table columns: %v", columns)

	// Check if last_name column exists
	var lastNameExists bool
	for _, column := range columns {
		if column == "last_name" {
			lastNameExists = true
			break
		}
	}

	if !lastNameExists {
		log.Println("Column last_name not found in users table. Adding...")

		// Add last_name column if it doesn't exist
		_, err = db.Exec("ALTER TABLE users ADD COLUMN last_name TEXT")
		if err != nil {
			log.Printf("Error adding last_name column: %v", err)
			return
		}

		// Update existing records
		_, err = db.Exec("UPDATE users SET last_name = '' WHERE last_name IS NULL")
		if err != nil {
			log.Printf("Error updating existing records: %v", err)
			return
		}

		log.Println("Column last_name successfully added to users table")
	} else {
		log.Println("Column last_name already exists in users table")
	}

	// Update test user data
	result, err := db.Exec("UPDATE users SET name = $1, last_name = $2 WHERE email = $3",
		"Тестовый", "Пользователь", "test@example.com")
	if err != nil {
		log.Printf("Error updating user data: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Updated users: %d", rowsAffected)
}

// CheckDatabase performs database checking and outputs information about tables
func CheckDatabase(db *sqlx.DB) {
	log.Println("Checking database...")

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	log.Println("Database connection successful")

	// Check users table
	var userCount int
	err := db.Get(&userCount, "SELECT COUNT(*) FROM users")
	if err != nil {
		log.Printf("Error checking users table: %v", err)
	} else {
		log.Printf("Number of users in database: %d", userCount)
	}

	// Check two_factor_codes table
	var codesCount int
	err = db.Get(&codesCount, "SELECT COUNT(*) FROM two_factor_codes")
	if err != nil {
		log.Printf("Error checking two_factor_codes table: %v", err)
	} else {
		log.Printf("Number of codes in database: %d", codesCount)
	}

	// Check categories table
	var categoriesCount int
	err = db.Get(&categoriesCount, "SELECT COUNT(*) FROM categories")
	if err != nil {
		log.Printf("Error checking categories table: %v", err)
	} else {
		log.Printf("Number of categories in database: %d", categoriesCount)
	}

	// Check purchases table
	var purchasesCount int
	err = db.Get(&purchasesCount, "SELECT COUNT(*) FROM purchases")
	if err != nil {
		log.Printf("Purchases table doesn't exist or other error: %v", err)
		// Check if status column exists in listings table
		var hasStatusColumn bool
		err = db.Get(&hasStatusColumn, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'listings' AND column_name = 'status'
			)
		`)
		if err != nil {
			log.Printf("Error checking status column in listings: %v", err)
		} else if !hasStatusColumn {
			log.Println("Need to run migration add_purchases.sql")
		}
	} else {
		log.Printf("Number of records in purchases table: %d", purchasesCount)
	}

	// Check users table structure
	var userColumns []string
	err = db.Select(&userColumns, `
		SELECT column_name FROM information_schema.columns 
		WHERE table_name = 'users'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Printf("Error checking users table structure: %v", err)
	} else {
		log.Printf("Users table structure: %v", userColumns)
	}

	// If no users exist, create a test user
	if userCount == 0 {
		createTestUser(db)
	}
}

// createTestUser creates a test user in the database
func createTestUser(db *sqlx.DB) {
	log.Println("Creating test user...")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	// Check if user exists
	var exists bool
	err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", "test@example.com")
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return
	}

	if exists {
		log.Println("Test user already exists")
		return
	}

	// Create user
	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, name, last_name, city, is_verified, created_at)
		VALUES ($1, $2, $3, $4, $5, true, NOW())
	`, "test@example.com", string(hashedPassword), "Тестовый", "Пользователь", "Москва")

	if err != nil {
		log.Printf("Error creating test user: %v", err)
		return
	}

	log.Println("Test user successfully created")

	// Output information about created user
	var userID int
	err = db.Get(&userID, "SELECT id FROM users WHERE email = $1", "test@example.com")
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		return
	}

	log.Printf("Test user ID: %d", userID)
}
