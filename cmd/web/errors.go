package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

func (app *App) ServerError(w http.ResponseWriter, err error) {
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *App) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *App) Unauthorized(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, []string{"error.page.html"}, &HTMLData{
		Title: "Unauthorized",
		Error: "Unauthorized",
	})
}

func (app *App) NotFound(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, []string{"error.page.html"}, &HTMLData{
		Title: "Page not found",
		Error: "Page not found",
	})
}

func (app *App) APINotFound(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(map[string]string{
		"error": "not found",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(j)
}

func (app *App) APIServerError(w http.ResponseWriter, err error) {
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	j, _ := json.Marshal(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(j)
}

func (app *App) APIClientError(w http.ResponseWriter, status int) {
	j, _ := json.Marshal(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    status,
			"message": http.StatusText(status),
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func (app *App) APIClientErrorWithMessage(w http.ResponseWriter, status int, message string) {
	j, _ := json.Marshal(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    status,
			"message": message,
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}
