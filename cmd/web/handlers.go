package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/code-mobi/board-checker/pkg/forms"
	"gitlab.com/code-mobi/board-checker/pkg/models"
)

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	worksheets, err := db.ListWorksheets()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	session := app.Sessions.Load(r)
	flash, err := session.PopString(w, "flash")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, []string{"home.page.html"}, &HTMLData{
		Flash:    flash,
		Worksheets: worksheets,
	})
}

func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	session := app.Sessions.Load(r)
	flash, err := session.PopString(w, "flash")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, []string{"login.page.html"}, &HTMLData{
		Flash: flash,
		Form:  &forms.LoginUser{},
	})
}

func (app *App) VerifyUser(w http.ResponseWriter, r *http.Request) {
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
		app.RenderHTML(w, r, []string{"login.page.html"}, &HTMLData{Form: form})
		return
	}

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	currentUserID, err := db.VerifyUser(form.Username, form.Password)
	if err == models.ErrInvalidCredentials {
		form.Failures["Generic"] = "Username or Password is incorrect"
		app.RenderHTML(w, r, []string{"login.page.html"}, &HTMLData{Form: form})
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	session := app.Sessions.Load(r)
	err = session.PutInt(w, "currentUserID", currentUserID)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) LogoutUser(w http.ResponseWriter, r *http.Request) {
	session := app.Sessions.Load(r)
	err := session.Remove(w, "currentUserID")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) ShowWorksheet(w http.ResponseWriter, r *http.Request) {
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	worksheet, err := db.GetWorksheet(worksheetID)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if worksheet == nil {
		app.NotFound(w, r)
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

	app.RenderHTML(w, r, []string{"worksheet.show.page.html", "worksheet.navbar.html", "photo.index.partial.html", "pagination.partial.html"}, &HTMLData{
		Worksheet: worksheet,
		Photos:  photos,
	})
}

func (app *App) EditWorksheet(w http.ResponseWriter, r *http.Request) {
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	if worksheetID == 0 {
		app.RenderHTML(w, r, []string{"worksheet.new.page.html", "worksheet.navbar.html"}, &HTMLData{})
		return
	}

	worksheet, err := db.GetWorksheet(worksheetID)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if worksheet == nil {
		app.NotFound(w, r)
		return
	}

	app.RenderHTML(w, r, []string{"worksheet.edit.page.html", "worksheet.navbar.html"}, &HTMLData{
		Worksheet: worksheet,
	})
}

func (app *App) SaveWorksheet(w http.ResponseWriter, r *http.Request) {
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.ServerError(w, err)
		return
	}

	decoder := form.NewDecoder()

	var f forms.Worksheet
	err := decoder.Decode(&f, r.PostForm)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	worksheet, err := db.GetWorksheet(worksheetID)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if worksheet == nil {
		worksheet = &models.Worksheet{
			ID:   f.ID,
			Name: f.Name,
		}
		err = db.InsertWorksheet(worksheet)
		if err != nil {
			app.ServerError(w, err)
			return
		}

	} else {
		worksheet = &models.Worksheet{
			ID:   worksheetID,
			Name: f.Name,
		}
		err = db.UpdateWorksheet(worksheet)
		if err != nil {
			app.ServerError(w, err)
			return
		}
	}

	session := app.Sessions.Load(r)
	err = session.PutString(w, "flash", "Worksheet was saved successfully!")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/worksheet/"+strconv.Itoa(worksheet.ID), http.StatusSeeOther)

}

func (app *App) NewPhoto(w http.ResponseWriter, r *http.Request) {
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	worksheet, err := db.GetWorksheet(worksheetID)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if worksheet == nil {
		app.NotFound(w, r)
		return
	}

	app.RenderHTML(w, r, []string{"photo.new.page.html", "worksheet.navbar.html"}, &HTMLData{
		Worksheet: worksheet,
	})
}

func (app *App) InsertPhoto(w http.ResponseWriter, r *http.Request) {
	session := app.Sessions.Load(r)
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

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
	if handler.Filename != "" {
		f, err := os.OpenFile(app.StoreDir+"/"+strconv.Itoa(worksheet.ID)+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		defer f.Close()
		io.Copy(f, uploadFile)
	} else {

		// form.Failures["Generic"] = "Please select file."
		// app.RenderHTML(w, r, []string{"photo.new.page.html"}, &HTMLData{Form: form})

		err = session.PutString(w, "flash", "Please choose file!")
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.RenderHTML(w, r, []string{"photo.new.page.html", "worksheet.navbar.html"}, &HTMLData{
			Worksheet: worksheet,
		})
		return
	}

	photo := &models.Photo{
		WorksheetID:     worksheet.ID,
		RunningNumber: runningNumber,
		FileName:      handler.Filename,
	}

	err = db.InsertPhoto(photo)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	err = session.PutString(w, "flash", "File was saved successfully!")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/worksheet/"+strconv.Itoa(worksheetID), http.StatusSeeOther)
}
