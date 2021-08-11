package main

import (
	"net/http"
)

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:  "Available",
		Env:     app.config.env,
		Version: version,
		Port:    app.config.port,
	}

	err := app.writeJSON(w, http.StatusOK, currentStatus, "status")

	if err != nil {
		app.errorJSON(w,err)
	}

}
