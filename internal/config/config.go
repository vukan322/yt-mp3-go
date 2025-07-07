package config

import "os"

type Config struct {
	Environment string
	Domain      string
	Port        string
	BasePath    string
}

func New() *Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "/yt-downloader"
	}

	return &Config{
		Environment: env,
		Domain:      domain,
		Port:        port,
		BasePath:    basePath,
	}
}
