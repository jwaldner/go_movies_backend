package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)


func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}


func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/v1/genres/:id", app.getAllMoviesByGenre)
	router.POST("/v1/admin/editmovie",app.wrap(secure.ThenFunc(app.editMovie)))
	router.POST("/v1/admin/editgenre",app.wrap(secure.ThenFunc(app.editGenre)))
	router.GET("/v1/admin/deletemovie/:id", app.wrap(secure.ThenFunc(app.deleteMovie)))

	router.GET("/v1/admin/movies",app.wrap(secure.ThenFunc(app.getAllMovies)))
	router.GET("/auth",app.wrap(secure.ThenFunc(app.statusHandler)))

	router.HandlerFunc(http.MethodPost, "/v1/graphql", app.moviesGraphQL)

	
	return app.enableCORS(router)
}
