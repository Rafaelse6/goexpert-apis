package database

import "github.com/Rafaelse6/goexpert/9-APIS/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByemail(email string) (*entity.User, error)
}
