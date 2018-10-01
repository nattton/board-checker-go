package main

import (
	"net/http"

	"gitlab.com/code-mobi/board-checker/pkg/models"
)

func (app *App) LoggedIn(r *http.Request) (bool, *models.User, error) {
	session := app.Sessions.Load(r)
	userID, err := session.GetInt("currentUserID")
	if err != nil {
		return false, nil, err
	}

	db := &models.Database{connect(app.DSN)}
	defer db.Close()

	user, err := db.UserInfo(userID)
	if err != nil {
		return false, nil, err
	}
	return true, user, nil
}

func (app *App) CurrentUser(r *http.Request) *models.User {
	if user := r.Context().Value(ctxUser); user != nil {
		return user.(*models.User)
	}
	return nil
}
