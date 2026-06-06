// Package user
package user

import (
	"context"
	"math/rand"
	"net/http"
	"time"
)

type contextKey string

const (
	userContextKey contextKey = "user"
	userCookieName string     = "user"
)

var generator *rand.Rand = rand.New(rand.NewSource(time.Now().Unix()))

func WithAuth(h http.Handler) http.Handler {
	userFn := func(w http.ResponseWriter, r *http.Request) {
		userCookie, err := r.Cookie(userCookieName)
		if err != nil {
			newUserCookie, err := CreateUserCookie(int(generator.Int31()))
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

		userContext := context.WithValue(r.Context(), userContextKey, userID)
		r = r.WithContext(userContext)

		userCookie.Expires = time.Now().Add(userAuthLifetime)

		http.SetCookie(w, userCookie)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(userFn)
}
