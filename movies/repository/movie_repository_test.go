package repository

import (
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

	userId := "456"
	movie := model.GetMovieDetailsResponse{
		Title:  "Inception",
		ImdbID: "tt1375666",
		Year:   "2010",
		Genre:  "Action, Sci-Fi",
		Actors: "Jackie Chan",
		Type:   "movie",
		Poster: "N/A",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO movies_cart (user_id, title, imdb_id, year, genre, actors, type, poster) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)).
		WithArgs(userId, movie.Title, movie.ImdbID, movie.Year, movie.Genre, movie.Actors, movie.Type, movie.Poster).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.AddToMovieCart(movie, userId)
	assert.NoError(t, err)
}

func TestAddToMovieCartDuplicate(t *testing.T) {
	db, mock, closeDb := setupMockDB(t)
	defer closeDb()

	repo := NewMovieRepository(db)

	userId := "456"
	movie := model.GetMovieDetailsResponse{
		Title:  "Inception",
		ImdbID: "tt1375666",
		Year:   "2010",
		Genre:  "Action, Sci-Fi",
		Actors: "Jackie Chan",
		Type:   "movie",
		Poster: "N/A",
	}

	pqErr := &pq.Error{Code: "23505"} // unique_violation
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO movies_cart (user_id, title, imdb_id, year, genre, actors, type, poster) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)).
		WithArgs(userId, movie.Title, movie.ImdbID, movie.Year, movie.Genre, movie.Actors, movie.Type, movie.Poster).
		WillReturnError(pqErr)

	err := repo.AddToMovieCart(movie, userId)
	assert.EqualError(t, err, "movie already added to the cart")
}

func TestGetMoviesInCartSuccess(t *testing.T) {
	db, mock, closeDb := setupMockDB(t)
	defer closeDb()

	repo := NewMovieRepository(db)

	movie := model.GetMovieDetailsResponse{
		Title:  "Inception",
		ImdbID: "tt1375666",
		Year:   "2010",
		Genre:  "Action, Sci-Fi",
		Actors: "Jackie Chan",
		Type:   "movie",
		Poster: "N/A",
	}
	userId := "123"
	rows := sqlmock.NewRows([]string{"title", "imdb_id", "year", "genre", "actors", "type", "poster"}).
		AddRow(
			movie.Title,
			movie.ImdbID,
			movie.Year,
			movie.Genre,
			movie.Actors,
			movie.Type,
			movie.Poster,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT title, imdb_id, year, genre, actors, type, poster FROM movies_cart WHERE user_id = $1")).
		WillReturnRows(rows)

	result, err := repo.GetMoviesInCart(userId)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, movie.Title, result[0].Title)
}
