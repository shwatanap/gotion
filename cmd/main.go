package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/shwatanap/gotion/internal/router"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Unable to read .env file: %v", err)
	}

	if err := router.Router().Run(":8000"); err != nil {
		panic(err)
	}
}
