package main

import (
	"log"
	"os"

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
	// initialisation de la base de données
	db, err := db.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

    // on ajoute Friendship et DirectMessage à la création !
	err = db.AutoMigrate(&models.User{}, &models.Friendship{}, &models.DirectMessage{}, &models.Match{})
    if err != nil {
        log.Fatal(err)
    }

	os.MkdirAll("uploads/avatars", os.ModePerm)
	// lancement des websockets
	hub := ws.NewHub(db)
	go hub.Run()

	// configuration de Gin
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
    corsConfig.AllowAllOrigins = true // permet à Vite d'y accéder
    corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(corsConfig))

	r.Static("/uploads", "./uploads")

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
		protected.POST("/match/save", social.SaveMatchHandler(db))
		protected.GET("/match/history", social.MyHistoryHandler(db))
		protected.GET("/match/leaderboard", social.LeaderboardHandler(db))
    }

	//r.Static("/", "./frontend/dist")

    log.Println("Serveur sur http://localhost:8081")
    log.Fatal(r.Run(":8081"))
}
