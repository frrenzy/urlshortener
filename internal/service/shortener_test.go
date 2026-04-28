package service

import (
	"context"
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
			got, gotErr := service.CreateShortURL(context.TODO(), test.original)

			require.NoError(t, gotErr, "should not fail")

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_shortenerService_GetOriginalURL(t *testing.T) {
	repository := repository.NewMapStorage()
	repository.Add(context.TODO(), model.NewLink(originalURL, shortURL))
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
			got, gotErr := service.GetOriginalURL(context.TODO(), test.short)

			if !test.wantErr {
				require.NoError(t, gotErr, "Should not fail")
			} else {
				require.Error(t, gotErr, "Should fail")
			}

			assert.Equal(t, test.want, got)
		})
	}
}
