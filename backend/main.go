package main

import (
	//"fmt"
	"log"
	//"net/http"

	"transcendance/internal/models"
	"transcendance/internal/ws"
	"transcendance/internal/db"
	"transcendance/internal/auth"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"

)

func main() {
//	// Notre route API qui ne peut pas être bloquée
//	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
//		// Ce message va s'afficher dans ton terminal de gauche !
//		fmt.Println("➡️ BINGO ! Requête reçue sur /api/hello depuis React !")
//
//		w.Header().Set("Content-Type", "application/json")
//		// On ajoute les autorisations de sécurité (CORS) au cas où
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		json.NewEncoder(w).Encode(map[string]string{"message": "Backend fonctionne"})
//	})
//
//	fs := http.FileServer(http.Dir("./frontend/dist"))
//	http.Handle("/", fs)
//
//	log.Println("Serveur sur http://localhost:8081")
//	log.Fatal(http.ListenAndServe(":8081", nil))

	//Database initialisation
	db, err := db.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	// auto-migrate models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	//run websockets
	hub := ws.NewHub()
	go hub.Run()

	//set up gin
	r := gin.Default()

	r.Use(cors.Default()) // change to set up cors properly
	
	r.POST("/signup", auth.SignupHandler(db))
	r.POST("/signin", auth.SigninHandler(db))
	r.Static("/", "./frontend/dist")

	log.Println("Serveur sur http://localhost:8081")
	log.Fatal(r.Run("localhost:8081"))
}
