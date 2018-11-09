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
		app.ClientError(w, err, http.StatusBadRequest)
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

type JSONWorksheets struct {
	models.Worksheets
	Host string
}

func (j JSONWorksheets) MarshalJSON() ([]byte, error) {
	type Worksheet struct {
		ID      int    `json:"id"`
		Number  string `json:"number"`
		Name    string `json:"name"`
		Created string `json:"created"`
	}
	worksheets := make([]Worksheet, len(j.Worksheets))
	for i, v := range j.Worksheets {
		worksheets[i] = Worksheet{
			ID:      v.ID,
			Number:  v.Number,
			Name:    v.Name,
			Created: v.Created.Format(time.RFC3339),
		}
	}
	return json.Marshal(worksheets)
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

func (app *App) APIListWorksheets(w http.ResponseWriter, r *http.Request) {
	db := &models.Database{connect(app.DSN)}
	defer db.Close()
	worksheets, err := db.ListWorksheets()
	if err == sql.ErrNoRows {
		app.APINotFound(w, r)
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"worksheets": JSONWorksheets{worksheets, "http://" + r.Host},
	})
	if err != nil {
		log.Fatal(err)
		app.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (app *App) APIListWorksheetsByTeam(w http.ResponseWriter, r *http.Request) {
	teamID, _ := strconv.Atoi(mux.Vars(r)["team_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()
	worksheets, err := db.ListWorksheetsByTeam(teamID)
	if err == sql.ErrNoRows {
		app.APINotFound(w, r)
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	b, err := json.Marshal(map[string]interface{}{
		"worksheets": JSONWorksheets{worksheets, "http://" + r.Host},
	})
	if err != nil {
		log.Fatal(err)
		app.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (app *App) APIShowWorksheet(w http.ResponseWriter, r *http.Request) {
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	worksheet, err := db.GetWorksheet(worksheetID)
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

	photos, err := db.ListPhotos(worksheet.ID, query)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	p := JSONPhotos{photos, "http://" + r.Host}
	b, err := json.Marshal(map[string]interface{}{
		"worksheet": worksheet,
		"photos":    p,
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
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	worksheet, err := db.GetWorksheet(worksheetID)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if worksheet == nil {
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

	fileDir := app.StoreDir + "/" + strconv.Itoa(worksheet.ID)
	os.MkdirAll(fileDir, os.ModePerm)
	f, err := os.OpenFile(fileDir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, uploadFile)

	photo := &models.Photo{
		WorksheetID:   worksheet.ID,
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
