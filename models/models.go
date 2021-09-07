package models

import (
	"database/sql"
	"time"
)

// Models is the wrapper for database
type Models struct {
	DB DBModel
}

// NewModels returns models with db pool
func NewModels(db *sql.DB) Models {

	return Models{
		DB: DBModel{DB: db},
	}

}

// Movie data model for a single movie
// Movie is the type for movies
type Movie struct {
	ID            int            `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Year          int            `json:"year"`
	ReleaseDate   time.Time      `json:"release_date"`
	Runtime       int            `json:"runtime"`
	Rating        int            `json:"rating"`
	MPAARating    string         `json:"mpaa_rating"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	MovieGenre    map[int]string `json:"genres"`
	MovieGenreIDs []int          `json:"genre_ids"`
	Poster        string         `json:"poster"`
}

// Genre data model for movie genres a lookup table
type Genre struct {
	ID        int       `json:"id"`
	GenreName string    `json:"genre_name"`
	JSONname  string    `json:"json_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MovieGenre data model for a movie genre
type MovieGenre struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	GenreID   int       `json:"genre_id"`
	Genre     Genre     `json:"genre"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User mock user data model
type User struct {
	ID       int
	Email    string
	Password string
}
