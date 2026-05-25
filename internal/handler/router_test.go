package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/service"
	"frrenzy/urlshortener/internal/util/user"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handlers_createShort(t *testing.T) {
	storage := repository.NewMapStorage()
	services := Services{
		URLService: service.NewShortenerService(storage),
	}
	router := NewRouter(services, user.WithAuth)
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
			errorBody:    errCanNotParseURL.Error(),
		},
		{
			name:         "wrong request body",
			body:         "123",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			errorBody:    errCanNotParseURL.Error(),
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
				assert.Equal(t, true, strings.HasPrefix(response.String(), config.Config.BaseAddress), "wrong response with short URL")
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
	storage.Add(t.Context(), model.NewLink(originalURL, shortURL, -1))
	services := Services{
		URLService: service.NewShortenerService(storage),
	}
	router := NewRouter(services, user.WithAuth)

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
			errorBody:    errNoShortURLFound.Error(),
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

func Test_handlers_createShortJSON(t *testing.T) {
	storage := repository.NewMapStorage()
	services := Services{
		URLService: service.NewShortenerService(storage),
	}
	router := NewRouter(services, user.WithAuth)
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
		body         any
		expectedCode int
		errorBody    string
		contentType  string
	}{
		{
			name:   "ok path",
			method: http.MethodPost,
			body: model.ShortenRequest{
				URL: "https://domain.com",
			},
			contentType:  "application/json",
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
			name:         "wrong Content-Type",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			errorBody:    errContentType.Error(),
		},
		{
			name:         "empty request body",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  "application/json",
			errorBody:    errCanNotDecode.Error(),
		},
		{
			name:         "not JSON request body",
			body:         "https://domain.com",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  "application/json",
			errorBody:    errCanNotDecode.Error(),
		},
		{
			name:         "wrong JSON request body",
			body:         `{ "url": "123" }`,
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  "application/json",
			errorBody:    errCanNotParseURL.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := client.R().SetBody(test.body)
			request.Method = test.method
			request.URL = "/api/shorten"
			if test.contentType != "" {
				request.Header.Set("Content-Type", test.contentType)
			}
			response, err := request.Send()

			require.NoError(t, err, "request error")

			assert.Equal(t, test.expectedCode, response.StatusCode())

			switch test.expectedCode {
			case http.StatusBadRequest:
				assert.Equal(t, test.errorBody, response.String())
			case http.StatusCreated:
				var resp model.ShortenResponse
				err = json.Unmarshal(response.Body(), &resp)

				require.NoError(t, err, "json parsing error")

				assert.Equal(t, true, strings.HasPrefix(resp.Result, config.Config.BaseAddress), "wrong response with short URL")
			}
		})
	}
}

func Test_handlers_createShortJSONBatch(t *testing.T) {
	storage := repository.NewMapStorage()
	services := Services{
		URLService: service.NewShortenerService(storage),
	}
	router := NewRouter(services, user.WithAuth)
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
		body         any
		expectedCode int
		errorBody    string
		contentType  string
	}{
		{
			name:   "ok path",
			method: http.MethodPost,
			body: model.BatchShortenRequest{
				{
					CorrelationID: "1",
					OriginalURL:   "https://domain1.com",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "https://domain2.com",
				},
			},
			contentType:  "application/json",
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
			name:         "wrong Content-Type",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			errorBody:    errContentType.Error(),
		},
		{
			name:         "empty request body",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  "application/json",
			errorBody:    errCanNotDecode.Error(),
		},
		{
			name:         "not JSON request body",
			body:         "https://domain.com",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  "application/json",
			errorBody:    errCanNotDecode.Error(),
		},
		{
			name:         "wrong JSON request body",
			body:         `{ "url": "123" }`,
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  "application/json",
			errorBody:    errCanNotDecode.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := client.R().SetBody(test.body)
			request.Method = test.method
			request.URL = "/api/shorten/batch"
			if test.contentType != "" {
				request.Header.Set("Content-Type", test.contentType)
			}
			response, err := request.Send()

			require.NoError(t, err, "request error")

			assert.Equal(t, test.expectedCode, response.StatusCode())

			switch test.expectedCode {
			case http.StatusBadRequest:
				assert.Equal(t, test.errorBody, response.String())
			case http.StatusCreated:
				var resp model.BatchShortenResponse
				err = json.Unmarshal(response.Body(), &resp)

				require.NoError(t, err, "json parsing error")

				assert.Equal(t, len(resp), 2, "should have 2 objects in array")

				body, ok := test.body.(model.BatchShortenRequest)
				require.Equal(t, true, ok)

				for i, res := range resp {
					assert.Equal(t, true, strings.HasPrefix(res.ShortURL, config.Config.BaseAddress), "wrong response with short URL")
					assert.Equal(t, body[i].CorrelationID, res.CorrelationID, "wrong response with short URL")
				}
			}
		})
	}
}

func Test_handlers_getAllUserURLs(t *testing.T) {
	repository := repository.NewMapStorage()
	mockLinks := []model.Link{
		model.NewLink(url.URL{Host: "domain.com"}, "wrong", 1),
		model.NewLink(url.URL{Host: "domain1.com"}, "right", 2),
		model.NewLink(url.URL{Host: "domain2.com"}, "also-right", 2),
	}
	for _, mockLink := range mockLinks {
		repository.Add(t.Context(), mockLink)
	}
	services := Services{
		URLService: service.NewShortenerService(repository),
	}
	router := NewRouter(services, user.WithAuth)

	server := httptest.NewServer(router)
	defer server.Close()

	client := resty.NewWithClient(&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}).SetBaseURL(server.URL).SetCookieJar(nil)

	tests := []struct {
		name         string
		method       string
		userID       int
		expectedCode int
		expectedData model.UserURLsResponce
	}{
		{
			name:         "ok path",
			method:       http.MethodGet,
			userID:       2,
			expectedCode: http.StatusOK,
			expectedData: model.UserURLsResponce{
				model.UserURLsResponceInstance{
					ShortURL:    "/" + mockLinks[1].ShortURL,
					OriginalURL: mockLinks[1].OriginalURL,
				},
				model.UserURLsResponceInstance{
					ShortURL:    "/" + mockLinks[2].ShortURL,
					OriginalURL: mockLinks[2].OriginalURL,
				},
			},
		},
		{
			name:         "method post",
			method:       http.MethodPost,
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
			name:         "user with no links",
			userID:       100,
			method:       http.MethodGet,
			expectedCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userCookie, err := user.CreateUserCookie(test.userID)
			require.NoError(t, err, "error creating cookie")
			request := client.R()
			request.Method = test.method
			request.URL = "/api/user/urls"
			request.SetCookie(userCookie)

			response, _ := request.Send()

			assert.Equal(t, test.expectedCode, response.StatusCode())

			if test.expectedCode == http.StatusOK {
				var resp model.UserURLsResponce
				err := json.Unmarshal(response.Body(), &resp)

				require.NoError(t, err, "json parsing error")
				assert.Equal(t, test.expectedData, resp)
			}
		})
	}
}

func Test_handlers_deleteAllUserURLs(t *testing.T) {
	repository := repository.NewMapStorage()
	mockLinks := []model.Link{
		model.NewLink(url.URL{Host: "domain.com"}, "wrong", 1),
		model.NewLink(url.URL{Host: "domain1.com"}, "right", 2),
		model.NewLink(url.URL{Host: "domain2.com"}, "also-right", 2),
	}
	for _, mockLink := range mockLinks {
		repository.Add(t.Context(), mockLink)
	}
	services := Services{
		URLService: service.NewShortenerService(repository),
	}
	router := NewRouter(services, user.WithAuth)

	server := httptest.NewServer(router)
	defer server.Close()

	client := resty.NewWithClient(&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}).SetBaseURL(server.URL).SetCookieJar(nil)

	tests := []struct {
		name         string
		userID       int
		data         string
		expectedCode int
		expectedData model.UserURLsResponce
	}{
		{
			name:         "ok path",
			userID:       2,
			data:         `["right"]`,
			expectedCode: http.StatusAccepted,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userCookie, err := user.CreateUserCookie(test.userID)
			require.NoError(t, err, "error creating cookie")
			request := client.R()
			request.Method = http.MethodDelete
			request.URL = "/api/user/urls"
			request.SetCookie(userCookie)
			request.Header.Set("Content-Type", "application/json")
			request.SetBody(test.data)

			response, _ := request.Send()

			assert.Equal(t, test.expectedCode, response.StatusCode())

			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			<-ticker.C

			deletedLink, err := repository.Get(t.Context(), "right")
			require.NoError(t, err, "should not fail")
			assert.Equal(t, true, deletedLink.DeletedFlag, "should be deleted")
		})
	}
}
