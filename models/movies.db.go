package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get gets a movie by id
func (m *DBModel) Get(id int) (*Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var movie Movie

	query := `
	SELECT id, title, description, "year", release_date, runtime, rating, mpaa_rating, created_at, updated_at, coalesce(poster, '')
	FROM movies WHERE id = $1;
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.Poster,
	)

	if err != nil {
		return &movie, err
	}

	genres, IDs, err := m.GetMovieGenres(id)

	if err != nil {
		return &movie, err
	}

	movie.MovieGenre = genres
	movie.MovieGenreIDs = IDs

	return &movie, nil
}

// All movies or one given an id
func (m *DBModel) All(genre ...int) ([]*Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""

	if len(genre) > 0 {
		where = fmt.Sprintf(`where id in (select movie_id from movies_genres where genre_id = %d)`, genre[0])
	}

	var movies []*Movie

	query := fmt.Sprintf(`
	SELECT id, title, description, "year", release_date, runtime, rating, mpaa_rating, created_at, updated_at, coalesce(poster,'') 
	FROM movies
	%s order by title;`, where)

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return movies, err
	}

	defer rows.Close()

	for rows.Next() {
		var i Movie

		err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Year,
			&i.ReleaseDate,
			&i.Runtime,
			&i.Rating,
			&i.MPAARating,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Poster,
		)

		if err != nil {
			return movies, err
		}

		genres, IDs, err := m.GetMovieGenres(i.ID)

		if err != nil {
			return movies, err
		}

		i.MovieGenre = genres
		i.MovieGenreIDs =IDs

		movies = append(movies, &i)
	}

	if err = rows.Err(); err != nil {
		return movies, err
	}

	return movies, nil
}

// GetMovieGenres get genres for a given movieID and potentially an error
func (m *DBModel) GetMovieGenres(movieID int) (map[int]string, []int,error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	genres := make(map[int]string)
	IDs := []int{}

	query := `
	SELECT movies_genres.id, movie_id, genre_id, g.genre_name
	,movies_genres.created_at,movies_genres.updated_at
	FROM movies_genres 
	LEFT JOIN genres g ON (g.id = movies_genres.genre_id)
	WHERE movie_id = $1;
	`

	rows, err := m.DB.QueryContext(ctx, query, movieID)

	defer rows.Close()

	if err != nil {
		return genres,IDs, err
	}

	for rows.Next() {
		var i MovieGenre

		err := rows.Scan(
			&i.ID,
			&i.MovieID,
			&i.GenreID,
			&i.Genre.GenreName,
			&i.CreatedAt,
			&i.UpdatedAt,
		)

		if err != nil {
			return genres, IDs, err
		}

		genres[i.ID] = i.Genre.GenreName
		 IDs = append(IDs, i.ID) 
	}

	if err = rows.Err(); err != nil {
		return genres, IDs, err
	}

	return genres, IDs, nil
}

// GenresAll All genres
func (m *DBModel) GenresAll() ([]*Genre, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var genres []*Genre

	query := `
	SELECT id, genre_name, json_name, created_at, updated_at
	FROM genres
	order by genre_name;
	`
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return genres, err
	}

	defer rows.Close()

	for rows.Next() {
		var i Genre

		err := rows.Scan(
			&i.ID,
			&i.GenreName,
			&i.JSONname,
			&i.CreatedAt,
			&i.UpdatedAt,
		)

		if err != nil {
			return genres, err
		}

		genres = append(genres, &i)
	}

	if err = rows.Err(); err != nil {
		return genres, err
	}

	return genres, nil
}

// InsertMovie inserts a single movie returns a potential error
func (m *DBModel) InsertMovie(movie Movie) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO movies
	(title, description, "year", release_date, runtime, rating, mpaa_rating
	,created_at, updated_at, poster)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);
	`
	_, err := m.DB.ExecContext(ctx, query,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.Poster,
	)

	if err != nil {
		return err
	}

	//genres, err := m.GetMovieGenres(id)

	//if err != nil {
	//	return movie, err
	//}

	//movie.MovieGenre = genres

	return nil
}

// UpdateMovie updates a single movie returns a potential error
func (m *DBModel) UpdateMovie(movie Movie, genres ...int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE movies
	SET title=$2, description=$3, "year"=$4, release_date=$5, runtime=$6, rating=$7, mpaa_rating=$8, updated_at=$9, poster=$10
	WHERE id=$1;
	`
	_, err := m.DB.ExecContext(ctx, query,
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		time.Now(),
		&movie.Poster,
	)

	if err != nil {
		return err
	}

	if len(genres) > 0 {
		err := m.ClearGenreFromMovie(movie.ID)

		if err != nil {
			return err
		}

		err = m.InsertMovieGenres(movie.ID, genres...)

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// DeleteMovie deletes a single movie given the ID, returns a potential error
func (m *DBModel) DeleteMovie(movieID int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM movies
	WHERE id=$1;
	`
	_, err := m.DB.ExecContext(ctx, query,
		movieID,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateGenres given a valid optional ID updates a genre otherwise
// inserts a single record, returns a potential error
func (m *DBModel) UpdateGenres(genre string, JSONname string, ID ...int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO genres
	(genre_name, json_name, created_at, updated_at)
	VALUES($1,$2,$3,$4)`

	if len(ID) > 0 {
		if ID[0] > 0 {
			query = `
	UPDATE genres SET genre_name=$2, jason_name=$3 updated_at=$4 WHERE id=$1;`

			_, err := m.DB.ExecContext(ctx, query,
				&ID[0],
				&genre,
				&JSONname,
				time.Now())

			if err != nil {
				return err
			} else {
				return nil
			}

		} else {
			return errors.New("invalid ID")
		}
	} else {
		_, err := m.DB.ExecContext(ctx, query,
			&genre,
			&JSONname,
			time.Now(),
			time.Now(),
		)

		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

// ClearGenreFromMovie clears the genres from a move, returns a potential error
func (m *DBModel) ClearGenreFromMovie(movieID int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM movies_genres
	WHERE movie_id=$1;
	`
	_, err := m.DB.ExecContext(ctx, query,
		movieID,
	)

	if err != nil {
		return err
	}

	return nil
}

// InsertMovieGenre  into a single movie returns a potential error
func (m *DBModel) InsertMovieGenre(movieID, genre int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO movies_genres
	(movie_id, genre_id, created_at, updated_at)
	VALUES($1,$2,$3,$4);
	`
	_, err := m.DB.ExecContext(ctx, query,
		&movieID,
		&genre,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// InsertMovieGenres genres into a single movie returns a potential error
func (m *DBModel) InsertMovieGenres(movieID int, genres ...int) error {

	if len(genres) > 0 {

		for _, g := range genres {
			err := m.InsertMovieGenre(movieID, g)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetGenre gets a genre by id
func (m *DBModel) GetGenre(id int) (*Genre, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var genre Genre

	query := `
	SELECT id, genre_name, json_name, created_at, updated_at
	FROM movies WHERE id = $1;
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&genre.ID,
		&genre.GenreName,
		&genre.JSONname,
		&genre.CreatedAt,
		&genre.UpdatedAt,
	)

	if err != nil {
		return &genre, err
	} else {
		return &genre, err
	}
}
