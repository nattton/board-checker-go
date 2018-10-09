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

	// Worksheet
	router.Handle("/worksheet/new", app.RequireLogin(http.HandlerFunc(app.NewWorksheet))).Methods("GET")
	router.Handle("/worksheet/new", app.RequireLogin(http.HandlerFunc(app.SaveWorksheet))).Methods("POST")

	worksheetRouter := router.PathPrefix("/worksheet/{worksheet_id:[0-9]+}").Subrouter()
	worksheetRouter.Handle("", app.RequireLogin(http.HandlerFunc(app.ShowWorksheet)))
	worksheetRouter.Handle("/edit", app.RequireLogin(http.HandlerFunc(app.EditWorksheet))).Methods("GET")
	worksheetRouter.Handle("/edit", app.RequireLogin(http.HandlerFunc(app.SaveWorksheet))).Methods("POST")
	worksheetRouter.Handle("/photo/new", app.RequireLogin(http.HandlerFunc(app.NewPhoto))).Methods("GET")
	worksheetRouter.Handle("/photo/new", app.RequireLogin(http.HandlerFunc(app.InsertPhoto))).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/user/login", app.APIUserLogin).Methods("POST")
	apiRouter.Handle("/worksheets", http.HandlerFunc(app.APIListWorksheets)).Methods("GET")
	apiRouter.Handle("/worksheet/{worksheet_id:[0-9]+}", http.HandlerFunc(app.APIShowWorksheet))
	apiRouter.Handle("/worksheet/{worksheet_id:[0-9]+}/photo/new", http.HandlerFunc(app.APIInsertPhoto)).Methods("POST")

	fileServer := http.FileServer(http.Dir(app.StaticDir))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	fileServer = http.FileServer(http.Dir(app.StoreDir))
	router.PathPrefix("/store/").Handler(http.StripPrefix("/store/", fileServer))

	router.NotFoundHandler = http.HandlerFunc(app.NotFound)

	return LogRequest(handlers.CompressHandler(SecureHeaders(app.LoggedInUser(router))))
}
