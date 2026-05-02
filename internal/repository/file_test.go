package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"testing"

	"frrenzy/urlshortener/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_fileStorage_Add(t *testing.T) {
	original := url.URL{
		Host: "domain.com",
	}
	short := "short-domain"
	expectedLink := model.NewLink(original, short)

	storageFile, err := os.CreateTemp(os.TempDir(), "storage.json")
	if err != nil {
		require.NoError(t, err, "could not create test storage file")
	}
	defer os.Remove(storageFile.Name())

	s := fileStorage{
		file: storageFile,
	}

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
			err := s.Add(context.TODO(), model.NewLink(test.original, test.short))
			require.NoError(t, err, "should not fail")

			storageFileContent, err := io.ReadAll(storageFile)
			require.NoError(t, err, "could not read test storage file")

			var links []model.Link
			err = json.Unmarshal(storageFileContent, &links)
			require.NoError(t, err, "could not Unmarshal test storage file")

			assert.Equal(t, 1, len(links), "should have 1 record")
			assert.Equal(t, expectedLink, links[0])
		})
	}
}

func Test_fileStorage_Get(t *testing.T) {
	mockLink := model.NewLink(url.URL{
		Host: "domain.com",
	}, "short-domain")

	storageFile, err := os.CreateTemp(os.TempDir(), "storage.json")
	if err != nil {
		require.NoError(t, err, "could not create test storage file")
	}
	defer os.Remove(storageFile.Name())

	storage, err := json.Marshal([]model.Link{mockLink})
	require.NoError(t, err, "could not Marshal test storage data")
	_, err = storageFile.Write(storage)
	require.NoError(t, err, "could not write test storage file")

	s := fileStorage{
		file: storageFile,
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
			got, gotErr := s.Get(context.TODO(), test.short)

			if !test.wantErr {
				require.NoError(t, gotErr, "Should not fail")
			} else {
				require.Error(t, gotErr, "Should fail")
			}

			assert.Equal(t, test.want, *got)
		})
	}
}

func Test_fileStorage_PingStorage(t *testing.T) {
	storageFile, err := os.CreateTemp(os.TempDir(), "storage.json")
	if err != nil {
		require.NoError(t, err, "could not create test storage file")
	}
	defer os.Remove(storageFile.Name())

	s := fileStorage{
		file: storageFile,
	}
	err = s.PingStorage(context.TODO())
	assert.ErrorIs(t, err, errDBNotConnected, "should fail")
}
