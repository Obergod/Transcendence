package main

import (
	"log"

	"transcendance/internal/auth"
	"transcendance/internal/db"
	"transcendance/internal/models"
	"transcendance/internal/social"
	"transcendance/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//"go.mongodb.org/mongo-driver/v2/x/mongo/driver/auth"
)

func main() {
	// Initialisation de la base de données
	db, err := db.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	// CORRECTION 1 : On ajoute Friendship et DirectMessage à la création !
	err = db.AutoMigrate(&models.User{}, &models.Friendship{}, &models.DirectMessage{})
	if err != nil {
		log.Fatal(err)
	}

	// Lancement des websockets
	hub := ws.NewHub()
	go hub.Run()

	// Configuration de Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Permet à Vite d'y accéder
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(corsConfig))

	// --- ROUTES PUBLIQUES ---
	r.POST("/api/signup", auth.SignupHandler(db))
	r.POST("/api/signin", auth.SigninHandler(db))

	// --- ROUTES PROTÉGÉES ---
	protected := r.Group("/api")

	r.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	protected.Use(auth.AuthRequired())
	{
		protected.PUT("/user/update", auth.UpdateProfileHandler(db))

		// CORRECTION 2 : On utilise les bons Handlers pour chaque route !
		protected.POST("/friends/request", social.SendFriendRequestHandler(db))
		protected.PUT("/friends/respond", social.RespondFriendRequestHandler(db))
		protected.GET("/friends/list", social.ListFriendsHandler(db))
		protected.GET("/user/me", auth.GetMeHandler(db))
	}

	//r.Static("/", "./frontend/dist")

	log.Println("Serveur sur http://localhost:8081")
	log.Fatal(r.Run(":8081"))
}