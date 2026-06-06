package user

import (
	"net/http"
	"time"
)

func CreateUserCookie(userID int) (*http.Cookie, error) {
	newUserJWT, err := buildJWTString(userID)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     userCookieName,
		Value:    newUserJWT,
		Expires:  time.Now().Add(userAuthLifetime),
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}, nil
}
