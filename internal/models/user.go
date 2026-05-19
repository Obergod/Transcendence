package models

type User struct {
	ID       int    `json:"id"`
	Pseudo   string `json:"pseudo"`
	Email    string `json:"email"`
	Password string `json:"password"` // Contiendra le hash Bcrypt
}

type UserRepository interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (User, error)
}