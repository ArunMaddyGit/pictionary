package main

import (
	"log"
	"net/http"

	"pictionary/handlers"
	"pictionary/store"
	"pictionary/ws"
)

func main() {
	log.Println("Pictionary backend starting...")
	ms := store.NewMemoryStore()
	var _ store.Store = ms

	hub := ws.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/join", handlers.HandleJoin(ms))
	mux.HandleFunc("/ws", ws.HandleWebSocket(hub))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
