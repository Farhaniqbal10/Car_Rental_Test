package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib" // Import driver pgx untuk sqlx
	"github.com/jmoiron/sqlx"
)

// Database struct menyimpan koneksi database sebagai *sqlx.DB
type Database struct {
	Conn *sqlx.DB
}

// NewDatabase - Konstruktor untuk membuat koneksi database
func NewDatabase() (*Database, error) {
	err := godotenv.Load("C:/belajar/car_rental_test/.env")
	if err != nil {
		log.Fatal("Error loading .env file from C:/belajar/car_rental_test/.env")
	}

	// Pastikan semua env variabel terisi
	requiredVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf("error: missing required environment variable: %s", v)
		}
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Gunakan sqlx untuk koneksi database
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)
	}

	fmt.Println("Database connected!")
	return &Database{Conn: db}, nil
}

// Close - Menutup koneksi database
func (db *Database) Close() {
	if db.Conn != nil {
		db.Conn.Close()
		fmt.Println("DB connection closed")
	}
}
