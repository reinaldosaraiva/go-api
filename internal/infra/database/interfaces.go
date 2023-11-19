package database

import (
	"github.com/reinaldosaraiva/go-api/internal/entity"
)
type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)

}