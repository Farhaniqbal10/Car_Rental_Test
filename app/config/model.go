package config

type DatabaseConfig struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type AppConfig struct {
	Mode     string `env:"APP_MODE"`
	Port     string `env:"APP_PORT"`
	Name     string `env:"APP_NAME"`
	Name_api string `env:"APP_NAME_API"`
	Url      string `env:"APP_URL"`
}

type Config struct {
	App AppConfig
	DB  DatabaseConfig
}
