package main

import (
	"log"
	"net/http"
)

func main() {
	// 1. Serveur minimal de secours pour faire patienter Docker
	// Tu le remplaceras par le vrai main.go de ton binôme lors du prochain "git pull"
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Backend temporaire en attente du vrai code de gbakulin"))
	})

	// 2. On ouvre le port 8081 pour que le docker-compose soit content
	log.Println("Serveur Go de secours démarré sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}