package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(rw, r)
	})
}

func (app application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Vary", "Authorization")
		
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// could set an anonymous user

		}

		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 {
			app.errorJSON(rw, errors.New("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			app.errorJSON(rw, errors.New("unauthorized - no bearer"))
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))

		if err != nil {
			app.errorJSON(rw, errors.New("unauthorized - failed HMACCheck"), http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			app.errorJSON(rw, errors.New("unauthorized - token expired"), http.StatusForbidden)
			app.logger.Println("unauthorized - token expired")
			return
		}

		
		if !claims.AcceptAudience("mydomain.com") {
			app.errorJSON(rw, errors.New("unauthorized - invalid audience"), http.StatusForbidden)
			return
		}

		if claims.Issuer != "mydomain.com" {
			app.errorJSON(rw, errors.New("unauthorized - invalid issuer"), http.StatusForbidden)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)

		if err != nil {
			app.errorJSON(rw, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		log.Println(userID)

		next.ServeHTTP(rw, r)
	})
}


func (app application) hasToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Vary", "Authorization")
		
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// could set an anonymous user

		}

		next.ServeHTTP(rw, r)
		
	})
}

