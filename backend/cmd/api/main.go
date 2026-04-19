package main

import (
	"log"

	"github.com/avaswani-build/fair-winds-api/internal/api"
	"github.com/avaswani-build/fair-winds-api/internal/weather"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	client := weather.NewStormglassClient()
	router := api.NewRouter(client)

	log.Println("server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
