package repository

import (
	"errors"
	"go-movie-api/movies/model"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MovieRespository interface {
	AddToMovieCart(movie model.GetMovieDetailsResponse, userId string) error
	GetMoviesInCart(userId string) (movies []model.MovieDetailsInCart, err error)
}

type movieRespository struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) movieRespository {
	return movieRespository{db: db}
}

func (mr movieRespository) AddToMovieCart(movie model.GetMovieDetailsResponse, userId string) error {
	_, err := mr.db.Exec(
		"INSERT INTO movies_cart (user_id, title, imdb_id, year, genre, actors, type, poster) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		userId, movie.Title, movie.ImdbID, movie.Year, movie.Genre, movie.Actors, movie.Type, movie.Poster,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return errors.New("movie already added to the cart")
		}
		log.Println(err)
		return err
	}

	return nil
}

func (mr movieRespository) GetMoviesInCart(userId string) (result []model.MovieDetailsInCart, err error) {
	rows, err := mr.db.Query(`SELECT title, imdb_id, year, genre, actors, type, poster FROM movies_cart WHERE user_id = $1`, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var movies []model.MovieDetailsInCart
	for rows.Next() {
		log.Println("Scan row:", err)
		var movie model.MovieDetailsInCart
		if err := rows.Scan(&movie.Title, &movie.ImdbID, &movie.Year, &movie.Genre, &movie.Type, &movie.Actors, &movie.Poster); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
