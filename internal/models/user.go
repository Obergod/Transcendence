package models

import (
	"gorm.io/gorm"
	"time"
)

type Friendship struct {
	ID			uint		`gorm:"primaryKey"`
	UserID		uint		`gorm:"index"`
	FriendID	uint		`gorm:"index"`
	Status		string		`gorm:"type:varchar(20);default:'pending'"`

	// Utilisation de la syntaxe complète et stricte pour GORM :
	User		*User		`gorm:"foreignKey:UserID;references:ID"`
	Friend		*User		`gorm:"foreignKey:FriendID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DirectMessage struct {
	ID			uint	`gorm:"primaryKey"`
	SenderID	uint	`gorm:"index"`
	ReceiverID	uint	`gorm:"index"`
	Content		string	`gorm:"type:text;not null"`
	CreatedAt	time.Time
}

type User struct {
	ID		 uint `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email	 string `gorm:"unique"`
	PasswordHash string `gorm:"not null"`
	AvatarURL	string
}

func CreateUser(db *gorm.DB, user *User) (error) {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByID(db *gorm.DB, userID uint) (*User, error) {
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByLogin(db *gorm.DB, login string) (*User, error) {
	var user User
	result := db.Where("email = ? OR username = ?", login, login).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *User) error {
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteUser(db *gorm.DB, user *User) error {
	result := db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
