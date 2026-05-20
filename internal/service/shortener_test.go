package service

import (
	"net/url"
	"testing"

	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	shortURL    = "short"
	originalURL = url.URL{
		Host: "domain.com",
	}
)

type mockGenerator struct{}

func (s mockGenerator) generateShort() string {
	return shortURL
}

func Test_shortenerService_CreateShortURL(t *testing.T) {
	repository := repository.NewMapStorage()
	service := ShortenerService{
		storage:   repository,
		generator: mockGenerator{},
	}

	tests := []struct {
		name     string
		original url.URL
		want     string
	}{
		{
			name: "simple test",
			original: url.URL{
				Host: "domain.com",
			},
			want: shortURL,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := service.CreateShortURL(t.Context(), test.original, -1)

			require.NoError(t, gotErr, "should not fail")

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_shortenerService_BatchCreateShortURL(t *testing.T) {
	repository := repository.NewMapStorage()
	service := ShortenerService{
		storage:   repository,
		generator: mockGenerator{},
	}

	firstOriginal := url.URL{
		Host: "domain.com",
	}
	firstID := "first"
	firstExpectedLink := model.NewLinkWithUUID(firstOriginal, shortURL, -1, firstID)
	secondOriginal := url.URL{
		Host: "another-domain.com",
	}
	secondID := "second"
	secondExpectedLink := model.NewLinkWithUUID(secondOriginal, shortURL, -1, secondID)
	expectedLinks := []model.Link{firstExpectedLink, secondExpectedLink}

	tests := []struct {
		name  string
		input []URLBatchInstance
		want  []model.Link
	}{
		{
			name: "simple test",
			input: []URLBatchInstance{
				{
					Original:      firstOriginal,
					CorrelationID: firstID,
				},
				{
					Original:      secondOriginal,
					CorrelationID: secondID,
				},
			},
			want: expectedLinks,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := service.BatchCreateShortURL(t.Context(), test.input, -1)

			require.NoError(t, gotErr, "should not fail")

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_shortenerService_GetOriginalURL(t *testing.T) {
	repository := repository.NewMapStorage()
	repository.Add(t.Context(), model.NewLink(originalURL, shortURL, -1))
	service := ShortenerService{
		storage:   repository,
		generator: mockGenerator{},
	}

	tests := []struct {
		name    string
		short   string
		want    string
		wantErr bool
	}{
		{
			name:    "simple test",
			short:   shortURL,
			want:    originalURL.String(),
			wantErr: false,
		},
		{
			name:    "nonexistent key",
			short:   "key",
			want:    "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := service.GetOriginalURL(t.Context(), test.short)

			if !test.wantErr {
				require.NoError(t, gotErr, "Should not fail")
			} else {
				require.Error(t, gotErr, "Should fail")
			}

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_shortenerService_GetByUser(t *testing.T) {
	repository := repository.NewMapStorage()
	mockLinks := []model.Link{
		model.NewLink(url.URL{Host: "domain.com"}, "wrong", 1),
		model.NewLink(url.URL{Host: "domain1.com"}, "right", 2),
		model.NewLink(url.URL{Host: "domain2.com"}, "also-right", 2),
	}
	for _, mockLink := range mockLinks {
		repository.Add(t.Context(), mockLink)
	}
	service := ShortenerService{
		storage:   repository,
		generator: mockGenerator{},
	}

	tests := []struct {
		name    string
		userID  int
		want    string
		wantErr bool
	}{
		{
			name:    "simple test",
			userID:  2,
			want:    originalURL.String(),
			wantErr: false,
		},
		{
			name:    "nonexistent key",
			userID:  10,
			want:    "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := service.GetByUser(t.Context(), test.userID)

			if !test.wantErr {
				require.NoError(t, gotErr, "Should not fail")
			} else {
				require.Error(t, gotErr, "Should fail")
			}

			assert.Equal(t, test.want, got)
		})
	}
}
