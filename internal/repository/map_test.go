package repository

import (
	"net/url"
	"testing"

	"frrenzy/urlshortener/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mapStorage_Add(t *testing.T) {
	original := url.URL{
		Host: "domain.com",
	}
	short := "short-domain"
	expectedLink := model.NewLink(original, short)

	tests := []struct {
		name     string
		original url.URL
		short    string
	}{
		{
			name:     "simple test",
			original: original,
			short:    short,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := NewMapStorage()
			s.Add(t.Context(), model.NewLink(test.original, test.short))

			assert.Equal(t, expectedLink, s.storage[test.short])
		})
	}
}

func Test_mapStorage_BatchAdd(t *testing.T) {
	firstOriginal := url.URL{
		Host: "domain.com",
	}
	firstShort := "short-domain"
	firstExpectedLink := model.NewLink(firstOriginal, firstShort)
	secondOriginal := url.URL{
		Host: "another-domain.com",
	}
	secondShort := "another-short-domain"
	secondExpectedLink := model.NewLink(secondOriginal, secondShort)
	links := []model.Link{firstExpectedLink, secondExpectedLink}

	tests := []struct {
		name  string
		links []model.Link
		short []string
	}{
		{
			name:  "simple test",
			links: links,
			short: []string{firstShort, secondShort},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := NewMapStorage()
			_, err := s.BatchAdd(t.Context(), links)
			require.NoError(t, err, "should not fail")

			assert.Equal(t, 2, len(s.storage), "should have 2 records")
			assert.Equal(t, firstExpectedLink, s.storage[test.short[0]])
			assert.Equal(t, secondExpectedLink, s.storage[test.short[1]])
		})
	}
}

func Test_mapStorage_Get(t *testing.T) {
	mockLink := model.NewLink(url.URL{
		Host: "domain.com",
	}, "short-domain")

	s := NewMapStorage()
	s.storage = map[string]model.Link{
		"short-domain": mockLink,
	}

	tests := []struct {
		name    string
		short   string
		want    model.Link
		wantErr bool
	}{
		{
			name:    "simple test",
			short:   "short-domain",
			want:    mockLink,
			wantErr: false,
		},
		{
			name:    "nonexistent key",
			short:   "key",
			want:    model.Link{},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := s.Get(t.Context(), test.short)

			if !test.wantErr {
				require.NoError(t, gotErr, "Should not fail")
			} else {
				require.Error(t, gotErr, "Should fail")
			}

			assert.Equal(t, test.want, *got)
		})
	}
}

func Test_mapStorage_PingStorage(t *testing.T) {
	s := NewMapStorage()
	err := s.PingStorage(t.Context())
	assert.ErrorIs(t, err, errDBNotConnected, "should fail")
}
