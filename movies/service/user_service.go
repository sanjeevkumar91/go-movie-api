package service

import (
	"go-movie-api/movies/model"
	"go-movie-api/movies/repository"
	"log"
)

type userService struct {
	repository repository.UserRespository
}

type UserService interface {
	CreateUser(req model.CreateUserRequest) (err error)
	GetUsers() (users []model.User, err error)
}

func NewUserService(repository repository.UserRespository) userService {
	return userService{repository: repository}
}

func (ms userService) CreateUser(req model.CreateUserRequest) (err error) {
	dbErr := ms.repository.CreateUser(req)
	if dbErr != nil {
		return dbErr
	}

	return nil
}

func (ms userService) GetUsers() (users []model.User, err error) {
	users, dbErr := ms.repository.GetUsers()
	if dbErr != nil {
		log.Println(dbErr)
		return nil, dbErr
	}

	return users, nil
}
