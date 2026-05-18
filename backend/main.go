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
	dsn := "host=postgres user=postgres password=goodPassword123 dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Paris"

	//db := repository.NewMemoryRepo()
	db := repository.NewPostgresRepo(dsn)
	authHandler := handlers.NewAuthHandler(db) //Maati: je donne la vraie db au handler

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