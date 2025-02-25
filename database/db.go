package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("erorr loading .env file")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)

	}

	fmt.Println("database connected! ")
	return conn, nil
}

func CloseDB(conn *pgx.Conn) {
	if conn != nil {
		conn.Close(context.Background())
		fmt.Println("DB connection closed")
	}
}
