package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadConfig(env string) error {
	// Coba cari file .env.<env> di direktori saat ini
	envFile := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envFile)
	if err == nil {
		log.Printf("Loaded environment configuration from %s", envFile)
		return nil
	}

	// Jika gagal, coba cari di direktori deployment
	rootDir := "/app" // Lokasi default di dalam container Docker
	envFilePath := filepath.Join(rootDir, envFile)
	err = godotenv.Load(envFilePath)
	if err != nil {
		return fmt.Errorf("error loading %s file: %v", envFilePath, err)
	}

	log.Printf("Loaded environment configuration from %s", envFilePath)
	return nil
}

func LoadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: Environment variable %s not found", key)
	}
	return value
}

func getRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no go.mod file found)")
		}
		dir = parent
	}
}
