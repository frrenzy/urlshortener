package repository

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/util/logger"

	"go.uber.org/zap"
)

type fileStorage struct {
	Repository
	file *os.File
}

func (s *fileStorage) Add(ctx context.Context, l model.Link) error {
	data, err := s.ReadAll()
	if err != nil {
		return err
	}

	data = append(data, l)

	return s.WriteAll(data)
}

func (s *fileStorage) BatchAdd(ctx context.Context, links []model.Link) ([]model.Link, error) {
	data, err := s.ReadAll()
	if err != nil {
		return nil, err
	}

	data = append(data, links...)

	return links, s.WriteAll(data)
}

func (s *fileStorage) Get(ctx context.Context, short string) (*model.Link, error) {
	data, err := s.ReadAll()
	if err != nil {
		return &model.Link{}, err
	}

	for _, link := range data {
		if link.ShortURL == short {
			return &link, nil
		}
	}

	return &model.Link{}, errNotFound
}

func (s *fileStorage) ReadAll() ([]model.Link, error) {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	fileContent, err := io.ReadAll(s.file)
	if err != nil {
		return nil, err
	}

	if len(fileContent) == 0 {
		return []model.Link{}, nil
	}

	var data []model.Link
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *fileStorage) WriteAll(data []model.Link) error {
	if len(data) == 0 {
		return nil
	}

	err := s.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = s.file.Seek(0, 0)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = s.file.Write(jsonData)
	if err != nil {
		return err
	}

	_, err = s.file.Seek(0, 0)
	return err
}

func (s *fileStorage) Close() {
	s.file.Close()
}

func (s *fileStorage) PingStorage(ctx context.Context) error {
	return errDBNotConnected
}

func NewFileStorage(path string) *fileStorage {
	p, err := filepath.Abs(path)
	if err != nil {
		logger.Log.Fatal("can not instanciate fileStorage", zap.Error(err))
	}
	file, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		logger.Log.Fatal("can not instanciate fileStorage", zap.Error(err))
	}

	return &fileStorage{
		file: file,
	}
}
