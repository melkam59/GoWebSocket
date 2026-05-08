package main

import (
	"net/http"
	"ws/internal/handlers"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/ws", handlers.WsEndpoint)
	return mux
}