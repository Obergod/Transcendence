package main

import (
	//"fmt"
	"os"
	// "log"
	"net/http"
	// "crypto/tls"
	"github.com/rs/zerolog"
    "github.com/rs/zerolog/log"

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
	// mux.HandleFunc("/api/register", authHandler.Register)
	// mux.HandleFunc("/api/login", authHandler.Login)
	mux.Handle("/api/register", metrics.TrackRequests(metrics.TrackActiveUsers(metrics.TrackVisitDuration(metrics.TrackUniqueVisitors(http.HandlerFunc(authHandler.Register))))))
    mux.Handle("/api/login", metrics.TrackRequests(metrics.TrackActiveUsers(metrics.TrackVisitDuration(metrics.TrackUniqueVisitors(http.HandlerFunc(authHandler.Login))))))
    mux.Handle("/metrics", promhttp.Handler())

	// 4. LA NOUVELLE ROUTE WEBSOCKET
	// mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	// On autorise toutes les origines pour le développement (CORS)
	// 	ws.ServeWs(hub, w, r)
	// })
	mux.Handle("/ws", metrics.TrackRequests(metrics.TrackActiveUsers(metrics.TrackVisitDuration(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ws.ServeWs(hub, w, r)
    })))))

	handler := metrics.TrackRequests(metrics.TrackUniqueVisitors(mux))

	logDir := "/app/logs"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        log.Fatal().Err(err).Msg("Failed to create log directory")
    }

	logFile, err := os.OpenFile(logDir+"/main.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to open log file")
    }
    defer logFile.Close()

    // 7. Configure zerolog to write to file
    // Log to both file and stdout
    logWriter := zerolog.MultiLevelWriter(logFile, os.Stdout)
    log.Logger = zerolog.New(logWriter).With().Timestamp().Logger()




	log.Info().Msg("Serveur sur https://localhost:8081")
	// log.Println("Serveur sur http://localhost:8081")
	// log.Fatal(http.ListenAndServe(":8081", metrics.TrackRequests(mux)))
	// log.Fatal(http.ListenAndServe(":8081", handler))
	err = http.ListenAndServeTLS(":8081", "/app/ssl/localhost+2.pem", "/app/ssl/localhost+2-key.pem", handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
	// log.Fatal(http.ListenAndServeTLS(":8081", "/app/ssl/localhost+2.pem", "/app/ssl/localhost+2-key.pem", handler))
}