package repository

import (
	"context"
	"database/sql"

	"frrenzy/urlshortener/internal/model"
)

type dbStorage struct {
	Repository
	db *sql.DB
}

func (s *dbStorage) Add(ctx context.Context, l model.Link) error {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO links
		    (uuid, short_url, original_url)
		VALUES
		    ($1, $2, $3)`, l.UUID, l.ShortURL, l.OriginalURL)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		return errDb
	}

	return nil
}

func (s *dbStorage) Get(ctx context.Context, short string) (*model.Link, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT uuid
		     , short_url
		     , original_url
		FROM links
		WHERE short = $1`, short)

	link := model.Link{}

	err := row.Scan(link.UUID, link.ShortURL, link.OriginalURL)
	if err != nil {
		return &model.Link{}, errDb
	}

	return &link, nil
}

func (s *dbStorage) Close() {
	s.db.Close()
}

func (s *dbStorage) PingStorage(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func NewDBStorage(db *sql.DB) *dbStorage {
	return &dbStorage{
		db: db,
	}
}
