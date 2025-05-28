package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT	string
	API_KEY string
}


func LoadConfig() Config{
	godotenv.Load()
	return Config{
		PORT: os.Getenv("PORT"),
		API_KEY: os.Getenv("API_KEY"),
	}
}
