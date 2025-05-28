package repository

import (
	"encoding/json"
	"errors"
	"go-movie-api/movies/model"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MovieRespository interface {
	AddToMovieCart(movie model.GetMovieDetailsResponse) error
	GetMoviesInCart() (movies []model.GetMovieDetailsResponse, err error)
}

type movieRespository struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) movieRespository {
	return movieRespository{db: db}
}

func (mr movieRespository) AddToMovieCart(movie model.GetMovieDetailsResponse) error {
	jsonData, jsonMarshalErr := json.Marshal(movie)
	if jsonMarshalErr != nil {
		log.Println("Error marshaling:", jsonMarshalErr)
		return jsonMarshalErr
	}

	_, err := mr.db.Exec(
		"INSERT INTO movies_cart (title, imdb_id, year, genre, movie_details) VALUES ($1, $2, $3, $4, $5)",
		movie.Title, movie.ImdbID, movie.Year, movie.Genre, jsonData,
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

func (mr movieRespository) GetMoviesInCart() (result []model.GetMovieDetailsResponse, err error) {
	var movies []model.GetMovieDetailsResponse
	rows, err := mr.db.Query("SELECT movie_details FROM movies_cart")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rawData []byte
		if err := rows.Scan(&rawData); err != nil {
			log.Println("Scan error:", err)
			continue
		}

		var movie model.GetMovieDetailsResponse
		if err := json.Unmarshal(rawData, &movie); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
