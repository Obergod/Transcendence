// backend/handlers/auth.go
package handlers

import (
	"encoding/json"
	"net/http"
	"../models" // Vérifie que ce chemin correspond à ton go.mod
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler contient notre interface (la fausse DB pour l'instant)
type AuthHandler struct {
	repo models.UserRepository
}

func NewAuthHandler(repo models.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// --- INSCRIPTION ---
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format invalide"})
		return
	}

	// Appel à notre fausse DB (qui va hacher le mot de passe et sauvegarder)
	err := h.repo.CreateUser(newUser)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Compte créé avec succès !"})
}

// --- CONNEXION ---
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds models.User
	json.NewDecoder(r.Body).Decode(&creds)

	// 1. On cherche l'utilisateur par son email dans la fausse DB
	user, err := h.repo.GetUserByEmail(creds.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email ou mot de passe incorrect"})
		return
	}

	// 2. On compare le mot de passe tapé en clair avec le hash stocké
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email ou mot de passe incorrect"})
		return
	}

	// 3. Succès ! (Ici on générera le vrai JWT plus tard)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Connexion réussie !",
		"pseudo":  user.Pseudo,
		"token":   "faux_token_temporaire_en_attendant_le_vrai_jwt",
	})
}