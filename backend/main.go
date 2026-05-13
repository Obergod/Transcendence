package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"transcendance/internal/models"
	"transcendance/internal/db"
)

func main() {
	// Notre route API qui ne peut pas être bloquée
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		// Ce message va s'afficher dans ton terminal de gauche !
		fmt.Println("➡️ BINGO ! Requête reçue sur /api/hello depuis React !")

		w.Header().Set("Content-Type", "application/json")
		// On ajoute les autorisations de sécurité (CORS) au cas où
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]string{"message": "Backend fonctionne"})
	})

	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	log.Println("Serveur sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

	//Database initialisation
	db, err := db.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
}
