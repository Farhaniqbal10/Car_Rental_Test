package env

import (
	"os"
)

func IsDevelopment() bool {
	return os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "develop"
}

func IsStaging() bool {
	return os.Getenv("APP_ENV") == "staging"
}

func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}
