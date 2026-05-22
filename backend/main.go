package main

import (
	//"fmt"
	"log"

	"transcendance/internal/models"
	"transcendance/internal/ws"
	"transcendance/internal/db"
	"transcendance/internal/auth"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
	jwt "github.com/appleboy/gin-jwt/v3"
	"gorm.io/gorm"

)

func registerRoute(r *gin.Engine, handle *jwt.GinJWTMiddleware, db *gorm.DB) {
	public := r.Group("/api")
	{
		public.POST("/signup", auth.SignupHandler(db))
		public.POST("/signin", handle.LoginHandler())
	}

	auth := r.Group
}

func main() {
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

	authMiddleware, err := jwt.New(auth.InitParams(db))
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	
	r.Static("/", "./frontend/dist")

	log.Println("Serveur sur http://localhost:8081")
	log.Fatal(r.Run("localhost:8081"))
}
