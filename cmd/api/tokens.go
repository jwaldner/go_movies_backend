package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jwaldner/go-movies-backend/models"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = models.User{
	ID:       10,
	Email:    "someuser@someemail.com",
	Password: "$2a$12$XCJRqR1XGaTVeC.K8l3hIeiME.jJB7O0gx995RRVA5OXXEu34Jkwm",
}

// Credentials used to test wether a user is allowed to login
type Credentials struct {
	UserName string `json:"email"`
	Password string `json:"password"`
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	// in a real app you would query for a valid user here
	// check the password hash stored in the database against what was supplied
	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))

	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}


	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	// was 24 hours, change to test for expire 
	//claims.Expires = jwt.NewNumericTime(time.Now().Add(1 * time.Minute))
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string {"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256,[]byte(app.config.jwt.secret))
	
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes),"response")




}
