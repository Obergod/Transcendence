package repository

import (
	"errors"
	"log"
	"sync"
	"transcendance/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type MemoryRepo struct {
	users  map[string]models.User
	nextID int
	mu     sync.Mutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		users:  make(map[string]models.User),
		nextID: 1,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (repo *MemoryRepo) CreateUser(u models.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[u.Email]; exists {
		return errors.New("cet email est déjà utilisé")
	}

	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		log.Println("Erreur Bcrypt:", err)
		return errors.New("erreur interne")
	}

	u.Password = hashedPassword
	u.ID = repo.nextID
	repo.nextID++
	repo.users[u.Email] = u

	log.Printf("Nouvel utilisateur : %s (ID: %d)", u.Pseudo, u.ID)
	return nil
}

func (repo *MemoryRepo) GetUserByEmail(email string) (models.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, exists := repo.users[email]
	if !exists {
		return models.User{}, errors.New("utilisateur introuvable")
	}
	return user, nil
}