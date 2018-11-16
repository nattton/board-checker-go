package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	dates, err := db.ListDistinctDate()
	if err != nil {
		app.ServerError(w, err)
		return
	}

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
		Flash:      flash,
		Dates:      dates,
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
		app.ClientError(w, err, http.StatusBadRequest)
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

func (app *App) IndexTeam(w http.ResponseWriter, r *http.Request) {
	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	teams, err := db.ListTeams()
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

	app.RenderHTML(w, r, []string{"team.index.page.html"}, &HTMLData{
		Flash: flash,
		Teams: teams,
	})
}

func (app *App) NewTeam(w http.ResponseWriter, r *http.Request) {

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	app.RenderHTML(w, r, []string{"team.new.page.html"}, &HTMLData{})
}

func (app *App) SaveTeam(w http.ResponseWriter, r *http.Request) {
	teamID, _ := strconv.Atoi(mux.Vars(r)["team_id"])

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

	var f forms.Team
	err := decoder.Decode(&f, r.PostForm)
	if err != nil {
		app.ClientError(w, err, http.StatusBadRequest)
		return
	}

	team, err := db.GetTeam(teamID)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if team == nil {
		team = &models.Team{
			Name: f.Name,
		}
		err = db.InsertTeam(team)
		if err != nil {
			app.ServerError(w, err)
			return
		}

	} else {
		team = &models.Team{
			ID:   teamID,
			Name: f.Name,
		}
		err = db.UpdateTeam(team)
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

	http.Redirect(w, r, "/teams", http.StatusSeeOther)
}

func (app *App) EditTeam(w http.ResponseWriter, r *http.Request) {
	teamID, _ := strconv.Atoi(mux.Vars(r)["team_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	team, err := db.GetTeam(teamID)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if team == nil {
		app.NotFound(w, r)
		return
	}

	app.RenderHTML(w, r, []string{"team.edit.page.html"}, &HTMLData{
		Team: team,
	})
}

func (app *App) IndexZone(w http.ResponseWriter, r *http.Request) {
	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	zones, err := db.ListZones()
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

	app.RenderHTML(w, r, []string{"zone.index.page.html"}, &HTMLData{
		Flash: flash,
		Zones: zones,
	})
}

func (app *App) NewZone(w http.ResponseWriter, r *http.Request) {

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	app.RenderHTML(w, r, []string{"zone.new.page.html"}, &HTMLData{})
}

func (app *App) SaveZone(w http.ResponseWriter, r *http.Request) {
	zoneID, _ := strconv.Atoi(mux.Vars(r)["zone_id"])

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

	var f forms.Zone
	err := decoder.Decode(&f, r.PostForm)
	if err != nil {
		app.ClientError(w, err, http.StatusBadRequest)
		return
	}

	zone, err := db.GetZone(zoneID)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if zone == nil {
		zone = &models.Zone{
			Name: f.Name,
		}
		err = db.InsertZone(zone)
		if err != nil {
			app.ServerError(w, err)
			return
		}

	} else {
		zone = &models.Zone{
			ID:   zoneID,
			Name: f.Name,
		}
		err = db.UpdateZone(zone)
		if err != nil {
			app.ServerError(w, err)
			return
		}
	}

	session := app.Sessions.Load(r)
	err = session.PutString(w, "flash", "Zone was saved successfully!")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/zones", http.StatusSeeOther)
}

func (app *App) EditZone(w http.ResponseWriter, r *http.Request) {
	zoneID, _ := strconv.Atoi(mux.Vars(r)["zone_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	zone, err := db.GetZone(zoneID)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if zone == nil {
		app.NotFound(w, r)
		return
	}

	app.RenderHTML(w, r, []string{"zone.edit.page.html"}, &HTMLData{
		Zone: zone,
	})
}

func (app *App) IndexWorksheetByDate(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	worksheets, err := db.ListWorksheetsByDate(date)
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

	app.RenderHTML(w, r, []string{"worksheet.index.page.html"}, &HTMLData{
		Flash:      flash,
		Worksheets: worksheets,
	})
}

func (app *App) IndexWorksheetByZone(w http.ResponseWriter, r *http.Request) {
	zoneID, _ := strconv.Atoi(mux.Vars(r)["zone_id"])

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	worksheets, err := db.ListWorksheetsByZone(zoneID)
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

	app.RenderHTML(w, r, []string{"worksheet.index.page.html"}, &HTMLData{
		Flash:      flash,
		Worksheets: worksheets,
	})
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

	app.RenderHTML(w, r, []string{"worksheet.show.page.html", "worksheet.navbar.html", "photo.index.partial.html", "pagination.partial.html"},
		&HTMLData{
			Worksheet: worksheet,
			Photos:    photos,
		})
}

func (app *App) NewWorksheet(w http.ResponseWriter, r *http.Request) {

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user := app.CurrentUser(r)
	if user == nil {
		app.Unauthorized(w, r)
		return
	}

	zones, _ := db.ListZones()
	teams, _ := db.ListTeams()

	app.RenderHTML(w, r, []string{"worksheet.new.page.html", "worksheet.navbar.html"},
		&HTMLData{
			Zones: zones,
			Teams: teams,
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
		app.NotFound(w, r)
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

	zones, _ := db.ListZones()
	teams, _ := db.ListTeams()

	app.RenderHTML(w, r, []string{"worksheet.edit.page.html", "worksheet.navbar.html"}, &HTMLData{
		Worksheet: worksheet,
		Zones:     zones,
		Teams:     teams,
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
		app.ClientError(w, err, http.StatusBadRequest)
		return
	}

	worksheet, err := db.GetWorksheet(worksheetID)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if worksheet == nil {
		worksheet = &models.Worksheet{
			Number: f.Number,
			Name:   f.Name,
			ZoneID: f.ZoneID,
			TeamID: f.TeamID,
		}
		err = db.InsertWorksheet(worksheet)
		if err != nil {
			app.ServerError(w, err)
			return
		}

	} else {
		worksheet = &models.Worksheet{
			ID:     worksheetID,
			Number: f.Number,
			Name:   f.Name,
			ZoneID: f.ZoneID,
			TeamID: f.TeamID,
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
		WorksheetID:   worksheet.ID,
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

func (app *App) DownloadPhoto(w http.ResponseWriter, r *http.Request) {
	worksheetID, _ := strconv.Atoi(mux.Vars(r)["worksheet_id"])

	var files []string

	path := app.StoreDir + "/" + strconv.Itoa(worksheetID)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}

	zipFileName := "photo_" + strconv.Itoa(worksheetID) + ".zip"
	downloadPath, err := ZipFiles(app, zipFileName, files)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, downloadPath, http.StatusSeeOther)
}

func ZipFiles(app *App, filename string, files []string) (string, error) {
	fileDir := app.StoreDir + "/temp"
	os.MkdirAll(fileDir, os.ModePerm)
	newZipFile, err := os.Create(fileDir + "/" + filename)
	if err != nil {
		return "", err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		zipfile, err := os.Open(file)
		if err != nil {
			return "", err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return "", err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return "", err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = info.Name()

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return "", err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return "", err
		}
	}
	return "/store/temp/" + filename, nil
}
