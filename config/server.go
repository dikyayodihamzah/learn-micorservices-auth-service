package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type serverConfig struct {
	Host string
	Port string
	URI  string
}

func NewServerConfig() serverConfig {
	return serverConfig{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
		URI:  fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
	}
}
