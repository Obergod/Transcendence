// internal/repository/postgres_repo.go
package repository

import (
	"errors"
	"log"
	"transcendance/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

// Initialise la connexion à PostgreSQL
func NewPostgresRepo(dsn string) *PostgresRepo {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Impossible de se connecter à PostgreSQL :", err)
	}

	// Crée la table "users" automatiquement si elle n'existe pas
	db.AutoMigrate(&models.User{})

	log.Println("✅ Connecté à PostgreSQL !")
	return &PostgresRepo{db: db}
}

// Hachage du mot de passe avec Bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// --- Création d'un compte (Register) ---
func (repo *PostgresRepo) CreateUser(u models.User) error {
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return errors.New("erreur interne de cryptographie")
	}
	u.Password = hashedPassword

	// GORM va tenter de créer l'utilisateur. S'il échoue (ex: email déjà pris), il renvoie une erreur.
	result := repo.db.Create(&u)
	if result.Error != nil {
		return errors.New("cet email est déjà utilisé")
	}

	return nil
}

// --- Connexion d'un compte (Login) ---
func (repo *PostgresRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, errors.New("utilisateur introuvable")
	}
	return user, nil
}