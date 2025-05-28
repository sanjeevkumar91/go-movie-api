package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDB() *sqlx.DB {

	dbConnectStr := "user=postgres password=sanjeev dbname=postgres host=localhost port=5432 sslmode=disable"
	var err error

	db, err = sqlx.Connect("postgres", dbConnectStr)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
		panic(err)
	}

	log.Println("Connected to PostgreSQL via sqlx.")
	return db
}
