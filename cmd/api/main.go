package main

import (
	"log"

	"github.com/avaswani-build/fair-winds-api/internal/api"
)

func main() {
	router := api.NewRouter()

	log.Println("server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
