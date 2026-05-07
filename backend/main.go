package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Backend fonctionne"})
	})
	log.Println("Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}