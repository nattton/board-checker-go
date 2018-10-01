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

	router.Handle("/project/new", app.RequireLogin(http.HandlerFunc(app.EditProject))).Methods("GET")
	router.Handle("/project/new", app.RequireLogin(http.HandlerFunc(app.SaveProject))).Methods("POST")

	projectRouter := router.PathPrefix("/project/{project_id:[0-9]+}").Subrouter()
	projectRouter.Handle("", app.RequireLogin(http.HandlerFunc(app.ShowProject)))
	projectRouter.Handle("/edit", app.RequireLogin(http.HandlerFunc(app.EditProject))).Methods("GET")
	projectRouter.Handle("/edit", app.RequireLogin(http.HandlerFunc(app.SaveProject))).Methods("POST")
	projectRouter.Handle("/photo/new", app.RequireLogin(http.HandlerFunc(app.NewPhoto))).Methods("GET")
	projectRouter.Handle("/photo/new", app.RequireLogin(http.HandlerFunc(app.InsertPhoto))).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/user/login", app.APIUserLogin).Methods("POST")
	apiRouter.Handle("/projects", http.HandlerFunc(app.APIListProjects)).Methods("GET")
	apiRouter.Handle("/project/{project_id:[0-9]+}", http.HandlerFunc(app.APIShowProject))
	apiRouter.Handle("/project/{project_id:[0-9]+}/photo/new", http.HandlerFunc(app.APIInsertPhoto)).Methods("POST")

	fileServer := http.FileServer(http.Dir(app.StaticDir))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	fileServer = http.FileServer(http.Dir(app.StoreDir))
	router.PathPrefix("/store/").Handler(http.StripPrefix("/store/", fileServer))

	router.NotFoundHandler = http.HandlerFunc(app.NotFound)

	return LogRequest(handlers.CompressHandler(SecureHeaders(app.LoggedInUser(router))))
}
