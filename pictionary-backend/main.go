package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"pictionary/game"
	"pictionary/handlers"
	"pictionary/store"
	"pictionary/ws"
)

func main() {
	log.Println("Pictionary backend starting...")
	ms := store.NewMemoryStore() // a
	hub := ws.NewHub()           // b
	engine := &game.GameEngine{  // c
		Store: ms,
		Hub:   hub,
	}
	router := &ws.MessageRouter{ // d
		Engine: engine,
		Store:  ms,
	}
	hub.Router = router // e
	hub.Engine = engine
	go hub.Run() // f

	mux := http.NewServeMux()
	mux.HandleFunc("/api/join", handlers.HandleJoin(ms, engine))
	mux.HandleFunc("/ws", ws.HandleWebSocket(hub, ms))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
