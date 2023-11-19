package entity

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	gorm.Model
}

func NewUser(name, email, password string) (*User,error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error while generating password hash %s", err)
		return nil,err
	}
	return &User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	},nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}