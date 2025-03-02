package config

import (
	"errors"
	"fmt"

	"github.com/golobby/env/v2"
	"github.com/joho/godotenv"
)

func Init() (Config, error) {
	conf := Config{}

	// Load environment variables dari .env dengan path eksplisit
	err := godotenv.Load("C:/belajar/car_rental_test/.env")
	if err != nil {
		fmt.Println("[WARNING] .env file not found, trying OS environment variables")
	}

	// Load environment dari OS
	err = env.Feed(&conf)
	if err != nil {
		return conf, errors.New("error loading environment variables from OS")
	}

	// Jika APP_PORT kosong, gunakan default 8081
	if conf.App.Port == "" {
		conf.App.Port = "8081"
	}

	// Debugging: Cek nilai yang terbaca
	fmt.Printf("[DEBUG] APP_PORT: %s\n", conf.App.Port)
	fmt.Printf("[DEBUG] DB_DRIVER: %s\n", conf.DB.Driver)
	fmt.Printf("[DEBUG] DB_USER: %s\n", conf.DB.User)
	fmt.Printf("[DEBUG] DB_PASSWORD: %s\n", conf.DB.Password)

	// Cek apakah DB_DRIVER kosong
	if conf.DB.Driver == "" {
		return conf, errors.New("DB_DRIVER is missing in .env or environment variables")
	}

	return conf, nil
}
