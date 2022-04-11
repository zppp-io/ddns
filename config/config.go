package config

import (
	"github.com/joho/godotenv"
	"os"
)

var configs map[string]string

func GetConfig(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return configs[key]
}

func InitEnvConfig() {
	godotenv.Load()
}
