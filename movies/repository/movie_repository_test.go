package repository

import (
	"encoding/json"
	"go-movie-api/movies/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql db, %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")
	return sqlxDB, mock, func() {
		sqlxDB.Close()
	}
}

func TestAddToMovieCartSuccess(t *testing.T) {
	db, mock, closeDb := setupMockDB(t)
	defer closeDb()

	repo := NewMovieRepository(db)

	movie := model.GetMovieDetailsResponse{
		Title:  "Inception",
		ImdbID: "tt1375666",
		Year:   "2010",
		Genre:  "Action, Sci-Fi",
	}

	jsonData, _ := json.Marshal(movie)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO movies_cart (title, imdb_id, year, genre, movie_details) VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(movie.Title, movie.ImdbID, movie.Year, movie.Genre, jsonData).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.AddToMovieCart(movie)
	assert.NoError(t, err)
}

func TestAddToMovieCartDuplicate(t *testing.T) {
	db, mock, closeDb := setupMockDB(t)
	defer closeDb()

	repo := NewMovieRepository(db)

	movie := model.GetMovieDetailsResponse{Title: "Inception", ImdbID: "tt1375666", Year: "2010", Genre: "Action, Sci-Fi"}
	jsonData, _ := json.Marshal(movie)

	pqErr := &pq.Error{Code: "23505"} // unique_violation
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO movies_cart (title, imdb_id, year, genre, movie_details) VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(movie.Title, movie.ImdbID, movie.Year, movie.Genre, jsonData).
		WillReturnError(pqErr)

	err := repo.AddToMovieCart(movie)
	assert.EqualError(t, err, "movie already added to the cart")
}

func TestGetMoviesInCartSuccess(t *testing.T) {
	db, mock, closeDb := setupMockDB(t)
	defer closeDb()

	repo := NewMovieRepository(db)

	movie := model.GetMovieDetailsResponse{Title: "Inception", ImdbID: "tt1375666", Year: "2010", Genre: "Action, Sci-Fi"}
	raw, _ := json.Marshal(movie)

	rows := sqlmock.NewRows([]string{"movie_details"}).
		AddRow(raw)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT movie_details FROM movies_cart")).
		WillReturnRows(rows)

	result, err := repo.GetMoviesInCart()
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, movie.Title, result[0].Title)
}
