package main

import (
	//"fmt"
	"log"
	"net/http"

	"transcendance/internal/handlers"
	"transcendance/internal/repository"
	"transcendance/internal/ws"
)

func main() {
	// 1. Initialisation de la DB et des Handlers
	db := repository.NewMemoryRepo()
	authHandler := handlers.NewAuthHandler(db)

	// 2. WEBSOCKETS
	hub := ws.NewHub()
	go hub.Run() // On lance le Hub en arrière-plan (goroutine)

	// 3. Déclaration des routes HTTP
	http.HandleFunc("/api/register", authHandler.Register)
	http.HandleFunc("/api/login", authHandler.Login)

	// 4. LA NOUVELLE ROUTE WEBSOCKET
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// On autorise toutes les origines pour le développement (CORS)
		ws.ServeWs(hub, w, r)
	})

	log.Println("Serveur sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}