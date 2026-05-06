package repository

import (
	"context"
	"database/sql"
	"net/url"

	"frrenzy/urlshortener/internal/model"
)

type dbStorage struct {
	Repository
	db *sql.DB
}

const insertQuery = `
  INSERT INTO links
	    (uuid, short_url, original_url)
	VALUES
	    ($1, $2, $3)
	ON CONFLICT DO NOTHING`

func (s *dbStorage) Add(ctx context.Context, l model.Link) error {
	res, err := s.db.ExecContext(ctx, insertQuery, l.UUID, l.ShortURL, l.OriginalURL)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		return ErrDBExisting
	}

	return nil
}

func (s *dbStorage) BatchAdd(ctx context.Context, links []model.Link) ([]model.Link, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	sttm, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return nil, err
	}
	defer sttm.Close()

	insertedRows := 0
	for _, link := range links {
		res, err := sttm.ExecContext(ctx, link.UUID, link.ShortURL, link.OriginalURL)
		if err != nil {
			return nil, err
		}

		inserted, _ := res.RowsAffected()
		insertedRows += int(inserted)
	}

	if insertedRows != len(links) {
		return nil, ErrDBExisting
	}

	tx.Commit()
	return links, nil
}

func (s *dbStorage) Get(ctx context.Context, short string) (*model.Link, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT uuid
		     , short_url
		     , original_url
		FROM links
		WHERE short_url = $1`, short)

	link := model.Link{}

	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL)
	if err != nil {
		return &model.Link{}, errDBNotFound
	}

	return &link, nil
}

func (s *dbStorage) GetByOriginal(ctx context.Context, original url.URL) (*model.Link, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT uuid
		     , short_url
		     , original_url
		FROM links
		WHERE original_url = $1`, model.URL{URL: original})

	link := model.Link{}

	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL)
	if err != nil {
		return &model.Link{}, errDBNotFound
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
