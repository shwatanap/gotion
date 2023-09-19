package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/shwatanap/gotion/internal/router"
)

func main() {
	if os.Getenv("ENV") == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Unable to read .env file: %v", err)
		}
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	if err := router.Router().Run(":" + port); err != nil {
		panic(err)
	}
}
