package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handlers_createShort(t *testing.T) {
	router := NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	client := resty.NewWithClient(&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}).SetBaseURL(server.URL)

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
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "method put",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "method patch",
			method:       http.MethodPatch,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "method delete",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
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
			request := client.R().SetBody(test.body)
			request.Method = test.method
			request.URL = "/"
			response, err := request.Send()

			require.NoError(t, err, "request error")

			assert.Equal(t, test.expectedCode, response.StatusCode())

			switch test.expectedCode {
			case http.StatusBadRequest:
				assert.Equal(t, test.errorBody, response.String())
			case http.StatusCreated:
				assert.Equal(t, true, strings.HasPrefix(response.String(), "http://localhost:8080/"), "wrong response with short URL")
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
	router := chi.NewRouter()
	router.Get(`/{id}`, handlers.redirectToOriginal)

	server := httptest.NewServer(router)
	defer server.Close()

	errorRedirect := errors.New("HTTP redirect blocked")
	client := resty.New().SetBaseURL(server.URL).SetRedirectPolicy(resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errorRedirect
	}))

	tests := []struct {
		name             string
		method           string
		shortID          string
		expectedCode     int
		expectedLocation string
		errorBody        string
	}{
		{
			name:             "ok path",
			method:           http.MethodGet,
			shortID:          shortURL,
			expectedCode:     http.StatusTemporaryRedirect,
			expectedLocation: originalURL.String(),
		},
		{
			name:         "method post",
			method:       http.MethodPost,
			shortID:      shortURL,
			expectedCode: http.StatusMethodNotAllowed,
			errorBody:    "wrong method",
		},
		{
			name:         "method put",
			method:       http.MethodPut,
			shortID:      shortURL,
			expectedCode: http.StatusMethodNotAllowed,
			errorBody:    "wrong method",
		},
		{
			name:         "method patch",
			method:       http.MethodPatch,
			shortID:      shortURL,
			expectedCode: http.StatusMethodNotAllowed,
			errorBody:    "wrong method",
		},
		{
			name:         "method delete",
			method:       http.MethodDelete,
			shortID:      shortURL,
			expectedCode: http.StatusMethodNotAllowed,
			errorBody:    "wrong method",
		},
		{
			name:         "nonexistent path",
			shortID:      "somepath",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			errorBody:    "no short URL found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := client.R()
			request.Method = test.method
			request.URL = "/" + test.shortID

			response, _ := request.Send()

			assert.Equal(t, test.expectedCode, response.StatusCode())

			switch test.expectedCode {
			case http.StatusBadRequest:
				assert.Equal(t, test.errorBody, response.String())
			case http.StatusTemporaryRedirect:
				location, err := response.RawResponse.Location()
				assert.NoError(t, err, "Location header should be present")
				assert.Equal(t, test.expectedLocation, location.String())
			}
		})
	}
}
