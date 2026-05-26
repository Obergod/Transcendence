package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToPostgreSQL() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	var db *gorm.DB
	var err error

	// Le backend va tenter de se connecter 5 fois, avec 2 secondes de pause
	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Connexion à la BDD établie !")
			return db, nil
		}

		log.Printf("⏳ Base de données non prête (essai %d/5). Attente de 2s...", i)
		time.Sleep(2 * time.Second)
	}

	// Si ça échoue vraiment après 10 secondes, on renvoie l'erreur fatale
	return nil, fmt.Errorf("impossible de se connecter à la BDD après 5 essais : %w", err)
}