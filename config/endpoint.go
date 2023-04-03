package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var EndpointPrefix = os.Getenv("ENDPOINT_PREFIX")
