package config

import (
    "github.com/joho/godotenv"
    "log"
)

func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println("Please create a .env file with the required environment variables.")
    }
}