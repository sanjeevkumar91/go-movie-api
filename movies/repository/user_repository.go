package repository

import (
	"go-movie-api/movies/model"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserRespository interface {
	CreateUser(user model.CreateUserRequest) error
	GetUsers() (users []model.User, err error)
}

type userRespository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) userRespository {
	return userRespository{db: db}
}

func (mr userRespository) CreateUser(user model.CreateUserRequest) error {
	_, err := mr.db.Exec(
		"INSERT INTO users (user_name, email, country) VALUES ($1, $2, $3)",
		user.Name, user.Email, user.Country,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (mr userRespository) GetUsers() (result []model.User, err error) {
	rows, err := mr.db.Query(`SELECT id, user_name, email, country, created_at, updated_at FROM users`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		log.Println("Scan row:", err)
		var user model.User
		if err := rows.Scan(&user.UserId, &user.Name, &user.Email, &user.Country, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
