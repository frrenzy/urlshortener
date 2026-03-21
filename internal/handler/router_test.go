package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/service"

	"github.com/stretchr/testify/assert"
)

func Test_handlers_createShort(t *testing.T) {
	handlers := handler{
		urlService: service.NewShortenerService(repository.NewMapStorage()),
	}

	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		errorBody    string
	}{
		{
			name:         "ok path",
			method:       http.MethodPost,
			body:         "https://domain.com",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "method get",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "method put",
			method:       http.MethodPut,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "method patch",
			method:       http.MethodPatch,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "method delete",
			method:       http.MethodDelete,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "empty request body",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			errorBody:    "can not parse request",
		},
		{
			name:         "wrong request body",
			body:         "123",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			errorBody:    "can not parse request",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			handlers.createShort(w, r)

			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusBadRequest {
				assert.Equal(t, test.errorBody, w.Body.String())
			} else {
				assert.Equal(t, true, strings.HasPrefix(w.Body.String(), "http://localhost:8080/"))
			}
		})
	}
}

func Test_handlers_redirectToOriginal(t *testing.T) {
	var (
		shortURL    = "short"
		originalURL = url.URL{
			Host:   "domain.com",
			Scheme: "https",
		}
	)

	storage := repository.NewMapStorage()
	storage.Add(originalURL, shortURL)
	handlers := handler{
		urlService: service.NewShortenerService(storage),
	}

	tests := []struct {
		name             string
		method           string
		shortId          string
		expectedCode     int
		expectedLocation string
		errorBody        string
	}{
		{
			name:             "ok path",
			method:           http.MethodGet,
			shortId:          shortURL,
			expectedCode:     http.StatusTemporaryRedirect,
			expectedLocation: originalURL.String(),
		},
		{
			name:         "method post",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "method put",
			method:       http.MethodPut,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "method patch",
			method:       http.MethodPatch,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "method delete",
			method:       http.MethodDelete,
			expectedCode: http.StatusBadRequest,
			errorBody:    "wrong method",
		},
		{
			name:         "empty path",
			shortId:      "",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			errorBody:    "no short URL provided",
		},
		{
			name:         "nonexistent path",
			shortId:      "somepath",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			errorBody:    "no short URL found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/{id}", nil)
			r.SetPathValue("id", test.shortId)
			w := httptest.NewRecorder()

			handlers.redirectToOriginal(w, r)

			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusBadRequest {
				assert.Equal(t, test.errorBody, w.Body.String())
			} else {
				assert.Equal(t, test.expectedLocation, w.Header().Get("Location"))
			}
		})
	}
}
