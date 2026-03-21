package repository

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mapStorage_Add(t *testing.T) {
	tests := []struct {
		name     string
		original url.URL
		short    string
	}{
		{
			name: "simple test",
			original: url.URL{
				Host: "domain.com",
			},
			short: "short-domain",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := NewMapStorage()
			s.Add(test.original, test.short)

			assert.Equal(t, test.original, s.storage[test.short])
		})
	}
}

func Test_mapStorage_Get(t *testing.T) {
	tests := []struct {
		name    string
		short   string
		want    url.URL
		wantErr bool
	}{
		{
			name:  "simple test",
			short: "short-domain",
			want: url.URL{
				Host: "domain.com",
			},
			wantErr: false,
		},
		{
			name:    "nonexistent key",
			short:   "key",
			want:    url.URL{},
			wantErr: true,
		},
	}

	s := NewMapStorage()
	s.storage = map[string]url.URL{
		"short-domain": {
			Host: "domain.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := s.Get(test.short)

			if !test.wantErr {
				require.NoError(t, gotErr, "Should not fail")
			} else {
				require.Error(t, gotErr, "Should fail")
			}

			assert.Equal(t, test.want, got)
		})
	}
}
