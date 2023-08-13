package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/shwatanap/gotion/internal/server"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Unable to read .env file: %v", err)
	}

	if err := server.Router().Run(":8000"); err != nil {
		panic(err)
	}
}
