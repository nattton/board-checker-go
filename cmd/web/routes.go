package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (app *App) Routes() http.Handler {
	router := mux.NewRouter()
	router.Handle("/", app.RequireLogin(http.HandlerFunc(app.Home))).Methods("GET")
	router.HandleFunc("/user/login", app.LoginUser).Methods("GET")
	router.HandleFunc("/user/login", app.VerifyUser).Methods("POST")
	router.Handle("/user/logout", app.RequireLogin(http.HandlerFunc(app.LogoutUser))).Methods("POST")

	// Team
	router.Handle("/teams",
		app.RequireLogin(http.HandlerFunc(app.IndexTeam))).Methods("GET")
	router.Handle("/team/new",
		app.RequireLogin(http.HandlerFunc(app.NewTeam))).Methods("GET")
	router.Handle("/team/new",
		app.RequireLogin(http.HandlerFunc(app.SaveTeam))).Methods("POST")
	router.Handle("/team/{team_id:[0-9]+}/edit",
		app.RequireLogin(http.HandlerFunc(app.EditTeam))).Methods("GET")
	router.Handle("/team/{team_id:[0-9]+}/edit",
		app.RequireLogin(http.HandlerFunc(app.SaveTeam))).Methods("POST")

	// Zone
	router.Handle("/zones",
		app.RequireLogin(http.HandlerFunc(app.IndexZone))).Methods("GET")
	router.Handle("/zone/new",
		app.RequireLogin(http.HandlerFunc(app.NewZone))).Methods("GET")
	router.Handle("/zone/new",
		app.RequireLogin(http.HandlerFunc(app.SaveZone))).Methods("POST")
	router.Handle("/zone/{zone_id:[0-9]+}/edit",
		app.RequireLogin(http.HandlerFunc(app.EditZone))).Methods("GET")
	router.Handle("/zone/{zone_id:[0-9]+}/edit",
		app.RequireLogin(http.HandlerFunc(app.SaveZone))).Methods("POST")

	// Worksheet
	router.Handle("/worksheet/new",
		app.RequireLogin(http.HandlerFunc(app.NewWorksheet))).Methods("GET")
	router.Handle("/worksheet/new",
		app.RequireLogin(http.HandlerFunc(app.SaveWorksheet))).Methods("POST")
	router.Handle("/worksheet/date/{date}",
		app.RequireLogin(http.HandlerFunc(app.IndexWorksheetByDate))).Methods("GET")
	router.Handle("/worksheet/zone/{zone_id:[0-9]+}",
		app.RequireLogin(http.HandlerFunc(app.IndexWorksheetByZone))).Methods("GET")

	worksheetRouter := router.PathPrefix("/worksheet/{worksheet_id:[0-9]+}").Subrouter()
	worksheetRouter.Handle("",
		app.RequireLogin(http.HandlerFunc(app.ShowWorksheet)))
	worksheetRouter.Handle("/download",
		app.RequireLogin(http.HandlerFunc(app.DownloadPhoto))).Methods("GET")
	worksheetRouter.Handle("/edit",
		app.RequireLogin(http.HandlerFunc(app.EditWorksheet))).Methods("GET")
	worksheetRouter.Handle("/edit",
		app.RequireLogin(http.HandlerFunc(app.SaveWorksheet))).Methods("POST")
	worksheetRouter.Handle("/delete",
		app.RequireLogin(http.HandlerFunc(app.DeleteWorksheet))).Methods("POST")
	worksheetRouter.Handle("/photo/new",
		app.RequireLogin(http.HandlerFunc(app.NewPhoto))).Methods("GET")
	worksheetRouter.Handle("/photo/new",
		app.RequireLogin(http.HandlerFunc(app.InsertPhoto))).Methods("POST")

	// API
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/user/login", app.APIUserLogin).Methods("POST")
	apiRouter.Handle("/worksheets", http.HandlerFunc(app.APIListWorksheets)).Methods("GET")
	apiRouter.Handle("/worksheet/{worksheet_id:[0-9]+}", http.HandlerFunc(app.APIShowWorksheet)).Methods("GET")
	apiRouter.Handle("/worksheet/{worksheet_id:[0-9]+}/photo/new", http.HandlerFunc(app.APIInsertPhoto)).Methods("POST")
	apiRouter.Handle("/team/{team_id:[0-9]+}/worksheets", http.HandlerFunc(app.APIListWorksheetsByTeam)).Methods("GET")

	// File Static
	fileServer := http.FileServer(http.Dir(app.StaticDir))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// file Store
	fileServer = http.FileServer(http.Dir(app.StoreDir))
	router.PathPrefix("/store/").Handler(http.StripPrefix("/store/", fileServer))

	router.NotFoundHandler = http.HandlerFunc(app.NotFound)

	return LogRequest(handlers.CompressHandler(SecureHeaders(app.LoggedInUser(router))))
}
