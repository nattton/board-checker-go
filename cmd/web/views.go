package main

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"gitlab.com/code-mobi/board-checker/pkg/models"
)

type HTMLData struct {
	Title        string
	User         *models.User
	LoggedIn     bool
	HiddenNavBar bool
	Flash        string
	Error        string
	Path         string
	Form         interface{}
	Dates        []string
	Team         *models.Team
	Teams        models.Teams
	Zone         *models.Zone
	Zones        models.Zones
	Worksheet    *models.Worksheet
	Worksheets   models.Worksheets
	Photos       models.Photos
	FormFields   models.FormFields
	PageInfo     *models.PageInfo
}

func (app *App) RenderHTML(w http.ResponseWriter, r *http.Request, pages []string, data *HTMLData) {
	if data == nil {
		data = &HTMLData{}
	}

	data.Path = r.URL.Path

	if user := app.CurrentUser(r); user != nil {
		data.LoggedIn = true
		data.User = user
	}

	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
	}
	for i := range pages {
		files = append(files, filepath.Join(app.HTMLDir, pages[i]))
	}

	fm := template.FuncMap{
		"humanDate":   humanDate,
		"timeString":  timeString,
		"humanNumber": humanNumber,
	}

	ts, err := template.New("").Funcs(fm).ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}

func humanDate(t time.Time) string {
	return t.Add(7 * time.Hour).Format("02 Jan 2006 at 15:04")
}

func timeString(ts string) string {
	timeForm := "2006-01-02T15:04:05Z"
	t, err := time.Parse(timeForm, ts)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func humanNumber(amount float64) string {
	return humanize.Commaf(amount)
}
