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
			newUserJWT, err := buildJWTString(currentID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userCookie = &http.Cookie{
				Name:     userCookieName,
				Value:    newUserJWT,
				Expires:  time.Now().Add(userAuthLifetime),
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				HttpOnly: true,
			}
		}

		userID := getUserID(userCookie.Value)

		userContext := context.WithValue(r.Context(), UserContextKey, userID)
		r = r.WithContext(userContext)

		userCookie.Expires = time.Now().Add(userAuthLifetime)

		http.SetCookie(w, userCookie)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(userFn)
}
