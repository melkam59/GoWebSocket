package main

import (
	"log"
	"net/http"
	"ws/internal/handlers"
)

func main() {
	mux := routes()
	log.Println("starting server on port 8080")

	// START THE PROCESSOR IN THE BACKGROUND!
	go handlers.ListenToWsChannel()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
