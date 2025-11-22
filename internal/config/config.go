package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string
	ServerPort string
	DSN        string
	CryptKey   string
}

func MustNewConfig() Config {
	_ = godotenv.Load()

	return Config{
		Env:        mustGetEnv("ENVIRONMENT"),
		ServerPort: mustGetEnv("HTTP_PORT"),
		DSN:        mustGetEnv("DATABASE_URL"),
		CryptKey:   mustGetEnv("CRYPT_KEY"),
	}
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	return value
}
