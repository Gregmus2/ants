package main

import (
	"ants/internal"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Dotenv loaded successfully")
}

func main() {
	internal.Serve()
}
