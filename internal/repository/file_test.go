package repository

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"slices"
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
	expectedLink := model.NewLink(original, short, -1)

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
			err := s.Add(t.Context(), model.NewLink(test.original, test.short, -1))
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

func Test_fileStorage_BatchAdd(t *testing.T) {
	firstOriginal := url.URL{
		Host: "domain.com",
	}
	firstShort := "short-domain"
	firstExpectedLink := model.NewLink(firstOriginal, firstShort, -1)
	secondOriginal := url.URL{
		Host: "another-domain.com",
	}
	secondShort := "another-short-domain"
	secondExpectedLink := model.NewLink(secondOriginal, secondShort, -1)
	links := []model.Link{firstExpectedLink, secondExpectedLink}

	storageFile, err := os.CreateTemp(os.TempDir(), "storage.json")
	if err != nil {
		require.NoError(t, err, "could not create test storage file")
	}
	defer os.Remove(storageFile.Name())

	s := fileStorage{
		file: storageFile,
	}

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
			_, err := s.BatchAdd(t.Context(), links)
			require.NoError(t, err, "should not fail")

			storageFileContent, err := io.ReadAll(storageFile)
			require.NoError(t, err, "could not read test storage file")

			var gotLinks []model.Link
			err = json.Unmarshal(storageFileContent, &gotLinks)
			require.NoError(t, err, "could not Unmarshal test storage file")

			assert.Equal(t, 2, len(links), "should have 2 records")
			assert.Equal(t, firstExpectedLink, gotLinks[0])
			assert.Equal(t, secondExpectedLink, gotLinks[1])
		})
	}
}

func Test_fileStorage_Get(t *testing.T) {
	mockLink := model.NewLink(url.URL{
		Host: "domain.com",
	}, "short-domain", -1)

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

func Test_fileStorage_GetByUser(t *testing.T) {
	mockLinks := []model.Link{
		model.NewLink(url.URL{Host: "domain.com"}, "wrong", 1),
		model.NewLink(url.URL{Host: "domain1.com"}, "right", 2),
		model.NewLink(url.URL{Host: "domain2.com"}, "also-right", 2),
	}

	storageFile, err := os.CreateTemp(os.TempDir(), "storage.json")
	if err != nil {
		require.NoError(t, err, "could not create test storage file")
	}
	defer os.Remove(storageFile.Name())

	storage, err := json.Marshal(mockLinks)
	require.NoError(t, err, "could not Marshal test storage data")
	_, err = storageFile.Write(storage)
	require.NoError(t, err, "could not write test storage file")

	s := fileStorage{
		file: storageFile,
	}

	tests := []struct {
		name    string
		userID  int
		want    []model.Link
		wantErr bool
	}{
		{
			name:   "simple test",
			userID: 2,
			want:   []model.Link{mockLinks[1], mockLinks[2]},
		},
		{
			name:   "nonexistent key",
			userID: 10,
			want:   []model.Link{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := s.GetByUser(t.Context(), test.userID)

			require.NoError(t, gotErr, "Should not fail")

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_fileStorage_DeleteByUser(t *testing.T) {
	mockLinks := []model.Link{
		model.NewLink(url.URL{Host: "domain.com"}, "wrong", 1),
		model.NewLink(url.URL{Host: "domain1.com"}, "right", 2),
		model.NewLink(url.URL{Host: "domain2.com"}, "also-right", 2),
	}

	storageFile, err := os.CreateTemp(os.TempDir(), "storage.json")
	if err != nil {
		require.NoError(t, err, "could not create test storage file")
	}
	defer os.Remove(storageFile.Name())

	storage, err := json.Marshal(mockLinks)
	require.NoError(t, err, "could not Marshal test storage data")
	_, err = storageFile.Write(storage)
	require.NoError(t, err, "could not write test storage file")

	s := fileStorage{
		file: storageFile,
	}

	tests := []struct {
		name        string
		userID      int
		wantDeleted []string
	}{
		{
			name:        "simple test",
			userID:      2,
			wantDeleted: []string{"right"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := s.DeleteByUser(t.Context(), test.userID, test.wantDeleted)

			require.NoError(t, gotErr, "Should not fail")

			data, err := s.readAll()
			require.NoError(t, err, "should not fail")
			for _, link := range data {
				if slices.Contains(test.wantDeleted, link.ShortURL) {
					assert.Equal(t, true, link.DeletedFlag, "should be deleted")
				} else {
					assert.Equal(t, false, link.DeletedFlag, "should be deleted")
				}
			}
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
	err = s.PingStorage(t.Context())
	assert.ErrorIs(t, err, errDBNotConnected, "should fail")
}
