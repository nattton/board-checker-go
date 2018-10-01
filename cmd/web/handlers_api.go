package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/code-mobi/board-checker/pkg/forms"
	"gitlab.com/code-mobi/board-checker/pkg/models"
)

type UserClaims struct {
	UserID int    `json:"uid"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

func (app *App) APIUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.LoginUser{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	}

	if !form.Valid() {
		app.APIClientError(w, http.StatusBadRequest)
		return
	}

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	currentUserID, err := db.VerifyUser(form.Username, form.Password)
	if err == models.ErrInvalidCredentials {
		app.APIClientErrorWithMessage(w, http.StatusBadRequest, "Email or Password is incorrect")
		return
	} else if err != nil {
		app.APIServerError(w, err)
		return
	}

	user, err := db.UserInfo(currentUserID)
	if err != nil {
		app.APIServerError(w, err)
		return
	}

	userClaims := UserClaims{
		user.ID,
		user.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	tokenString, err := token.SignedString([]byte(app.SecretKey))
	if err != nil {
		app.ServerError(w, err)
		return
	}

	b, _ := json.Marshal(map[string]interface{}{
		"access_token": tokenString,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

type JSONProjects struct {
	models.Projects
	Host string
}

func (j JSONProjects) MarshalJSON() ([]byte, error) {
	type Project struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		FileURL string `json:"fileURL"`
		Created string `json:"created"`
	}
	projects := make([]Project, len(j.Projects))
	for i, v := range j.Projects {
		projects[i] = Project{
			ID:      v.ID,
			Name:    v.Name,
			FileURL: j.Host + v.FilePath(),
			Created: v.Created.Format(time.RFC3339),
		}
	}
	return json.Marshal(projects)
}

type JSONPhotos struct {
	models.Photos
	Host string
}

func (j JSONPhotos) MarshalJSON() ([]byte, error) {
	type Photo struct {
		ID            int    `json:"id"`
		RunningNumber int    `json:"runningNumber"`
		FileURL       string `json:"fileURL"`
		Created       string `json:"created"`
	}
	photos := make([]Photo, len(j.Photos))
	for i, v := range j.Photos {
		photos[i] = Photo{
			ID:            v.ID,
			RunningNumber: v.RunningNumber,
			FileURL:       j.Host + v.FilePath(),
			Created:       v.Created.Format(time.RFC3339),
		}
	}
	return json.Marshal(photos)
}

func (app *App) APIListProjects(w http.ResponseWriter, r *http.Request) {
	db := &models.Database{connect(app.DSN)}
	defer db.Close()
	log.Printf("Test")
	projects, err := db.ListProjects()
	if err == sql.ErrNoRows {
		app.APINotFound(w, r)
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"projects": JSONProjects{projects, "http://" + r.Host},
	})
	if err != nil {
		log.Fatal(err)
		app.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (app *App) APIShowProject(w http.ResponseWriter, r *http.Request) {
	projectID, _ := strconv.Atoi(mux.Vars(r)["project_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	project, err := db.GetProject(projectID)
	if err == sql.ErrNoRows {
		app.APINotFound(w, r)
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	query := forms.NewQuery()
	query.Q = r.FormValue("q")
	query.Start, _ = strconv.Atoi(r.FormValue("start"))
	maxResults, err := strconv.Atoi(r.FormValue("maxResults"))
	if err == nil {
		query.MaxResults = maxResults
	}

	photos, err := db.ListPhotos(project.ID, query)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	p := JSONPhotos{photos, "http://" + r.Host}
	b, err := json.Marshal(map[string]interface{}{
		"project": project,
		"photos":  p,
	})

	if err != nil {
		log.Fatal(err)
		app.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (app *App) APIInsertPhoto(w http.ResponseWriter, r *http.Request) {
	projectID, _ := strconv.Atoi(mux.Vars(r)["project_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	project, err := db.GetProject(projectID)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if project == nil {
		app.NotFound(w, r)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		app.ServerError(w, err)
		return
	}

	runningNumber, _ := strconv.Atoi(r.FormValue("running_number"))

	uploadFile, handler, err := r.FormFile("uploadFile")
	if err != nil {
		app.ServerError(w, err)
		return
	}
	defer uploadFile.Close()

	log.Printf("%v", handler.Header)

	fileDir := app.StoreDir + "/" + strconv.Itoa(project.ID)
	os.MkdirAll(fileDir, os.ModePerm)
	f, err := os.OpenFile(fileDir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, uploadFile)

	photo := &models.Photo{
		ProjectID:     project.ID,
		RunningNumber: runningNumber,
		FileName:      handler.Filename,
	}

	err = db.InsertPhoto(photo)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"status": "Success",
	})

	if err != nil {
		log.Fatal(err)
		app.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
