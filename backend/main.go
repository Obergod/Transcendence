package main

import (
	//"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"transcendance/internal/handlers"
	"transcendance/internal/repository"
	"transcendance/internal/ws"
	"transcendance/internal/metrics"
)

func main() {
	// 1. Initialisation de la DB et des Handlers
	db := repository.NewMemoryRepo()
	authHandler := handlers.NewAuthHandler(db)

	// 2. WEBSOCKETS
	hub := ws.NewHub()
	go hub.Run() // On lance le Hub en arrière-plan (goroutine)

	mux := http.NewServeMux()

	// 3. Déclaration des routes HTTP
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.Handle("/metrics", promhttp.Handler())

	// 4. LA NOUVELLE ROUTE WEBSOCKET
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// On autorise toutes les origines pour le développement (CORS)
		ws.ServeWs(hub, w, r)
	})

	log.Println("Serveur sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", metrics.TrackRequests(mux)))
}