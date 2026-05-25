// Package user
package user

import (
	"context"
	"net/http"
	"time"
)

type ContextKey string

const (
	UserContextKey ContextKey = "user"
	userCookieName string     = "user"
)

var currentID int = 1

func WithAuth(h http.Handler) http.Handler {
	userFn := func(w http.ResponseWriter, r *http.Request) {
		userCookie, err := r.Cookie(userCookieName)
		if err != nil {
			currentID += 1
			newUserCookie, err := CreateUserCookie(currentID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userCookie = newUserCookie
		}

		userID := getUserID(userCookie.Value)
		if userID == -1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userContext := context.WithValue(r.Context(), UserContextKey, userID)
		r = r.WithContext(userContext)

		userCookie.Expires = time.Now().Add(userAuthLifetime)

		http.SetCookie(w, userCookie)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(userFn)
}
