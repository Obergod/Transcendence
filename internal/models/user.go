package models

import (
)

type User struct {
	ID		 uint `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email	 string
}
