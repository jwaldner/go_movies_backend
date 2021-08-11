package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/jwaldner/go-movies-backend/models"
)

type JSONresp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type MoviePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Print(errors.New("invalid id parameter"))
		return
	}

	app.logger.Println("id:", id)

	movie, err := app.models.DB.Get(id)

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movie, "movie")

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {

	movies, err := app.models.DB.All()

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, "movies")

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
		return
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Print(errors.New("invalid id parameter"))
		return
	}

	err = app.models.DB.DeleteMovie(id)

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
		return
	}

	ok := JSONresp{OK: true, Message: "movie deleted"}

	err = app.writeJSON(w, http.StatusOK, ok, "response")

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Print(errors.New(err.Error()))
		return
	}
}

func (app *application) editMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload

	app.logger.Print("edit movie called")

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var movie models.Movie

	if payload.ID != "0" {
		id, _ := strconv.Atoi(payload.ID)
		m, _ := app.models.DB.Get(id)
		movie = *m
		movie.UpdatedAt = time.Now()
	}

	movie.ID, _ = strconv.Atoi(payload.ID)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, _ = strconv.Atoi(payload.Runtime)
	movie.Rating, _ = strconv.Atoi(payload.Rating)
	movie.MPAARating = payload.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	if movie.Poster == "" {
		movie = getPoster(movie)
	}

	if movie.ID == 0 {
		err = app.models.DB.InsertMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			app.logger.Println(err)
			return
		}
	} else {
		err = app.models.DB.UpdateMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			app.logger.Println(err)
			return
		}
	}

	ok := JSONresp{
		OK:      true,
		Message: "Movie saved",
	}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
		return
	}

	app.logger.Printf("movie ID: %v saved", movie.ID)
}

func (app *application) searchMovies(w http.ResponseWriter, r *http.Request) {
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {

	genres, err := app.models.DB.GenresAll()

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, genres, "genres")

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
	}
}

func (app *application) getAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Print(errors.New("invalid id parameter"))
		return
	}

	movies, err := app.models.DB.All(id)

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, "movies")

	if err != nil {
		app.errorJSON(w, err)
		app.logger.Println(err)
	}
}

func getPoster(movie models.Movie) models.Movie {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     string  `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			Popularity       float64 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			ReleaseDate      string  `json:"release_date"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
		} `json:"results"`
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
	}

	client := &http.Client{}
	key := os.Getenv("GO_MOVIES_DB_KEY")
	theURL := "https://api.themoviedb.org/3/search/movie?api_key="

	// just in case the url is bad
	//log.Println(theURL + key + "&query=" + url.QueryEscape(movie.Title))

	req, err := http.NewRequest("GET", theURL+key+"&query="+url.QueryEscape(movie.Title), nil)
	if err != nil {
		log.Println(err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}

	var responseObject TheMovieDB

	json.Unmarshal(bodyBytes, &responseObject)

	if len(responseObject.Results) > 0 {
		movie.Poster = responseObject.Results[0].PosterPath
	}

	return movie
}
