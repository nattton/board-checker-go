package main

import (
	"context"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"gitlab.com/code-mobi/board-checker/pkg/models"
)

type AppContext int

const (
	ctxUser AppContext = 1 + iota
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"RemoteAddr": r.RemoteAddr,
			"Proto":      r.Proto,
			"Method":     r.Method,
			"Request":    r.URL.RequestURI(),
		}).Info("LogRequest")
		// pattern := `%s - "%s %s %s"`
		// log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (app *App) LoggedInUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, user, err := app.LoggedIn(r)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUser, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *App) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.CurrentUser(r)
		if user == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString != "" {
			fmt.Sscanf(tokenString, "Bearer %s", &tokenString)
			token, _ := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.SecretKey), nil
			})

			if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {

				fmt.Printf("%v %v\n", claims.Name, claims.StandardClaims.ExpiresAt)

				ctx := context.WithValue(r.Context(), ctxUser, models.User{
					ID:   claims.UserID,
					Name: claims.Name,
				})
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			app.APIClientErrorWithMessage(w, http.StatusUnauthorized, "Authorization Invalid!")
			return
		}

		next.ServeHTTP(w, r)
	})
}
