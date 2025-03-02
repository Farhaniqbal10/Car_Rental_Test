package main

import (
	"fmt"
	"log"
	"time"

	"car_rental_test/app/config"
	"car_rental_test/app/driver"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Inisialisasi konfigurasi
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("[BOOT] Failed to initiate config | %s", err.Error())
	}

	// Inisialisasi router
	r := gin.Default()
	r.RedirectTrailingSlash = false

	// Menentukan mode aplikasi (debug/release)
	gin.SetMode(gin.ReleaseMode)
	if cfg.App.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	// Inisialisasi database
	dbCfg := cfg.DB
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name)

	db, err := sqlx.Open(dbCfg.Driver, dsn)
	if err != nil {
		log.Fatalf("[DATABASE] Error while opening connection to database | %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("[DATABASE] Error while checking connection to database | %s", err.Error())
	}

	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Minute)

	// Inisialisasi driver untuk menghubungkan repository, service, dan handler
	if err := driver.Init(r, db); err != nil {
		log.Fatalf("[BOOT] Failed to init modules driver | %s", err.Error())
	}

	fmt.Println("[DEBUG] Menggunakan database config:")
	fmt.Printf("[DEBUG] Driver: %s\n", dbCfg.Driver)
	fmt.Printf("[DEBUG] Host: %s\n", dbCfg.Host)
	fmt.Printf("[DEBUG] Port: %s\n", dbCfg.Port)
	fmt.Printf("[DEBUG] User: %s\n", dbCfg.User)
	fmt.Printf("[DEBUG] Password: %s\n", dbCfg.Password)
	fmt.Printf("[DEBUG] Name: %s\n", dbCfg.Name)

	// Menjalankan server
	port := cfg.App.Port
	log.Printf("[SERVER] Running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[SERVER] Failed to start | %s", err.Error())
	}
}
