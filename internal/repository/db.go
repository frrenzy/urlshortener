package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"frrenzy/urlshortener/internal/model"
)

type dbStorage struct {
	Repository
	db *sql.DB
}

func (s *dbStorage) Add(ctx context.Context, l model.Link) error {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO links
		    (uuid, short_url, original_url, user_id)
		VALUES
		    ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING`,
		l.UUID, l.ShortURL, l.OriginalURL, l.UserID,
	)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		return ErrDBExisting
	}

	return nil
}

func (s *dbStorage) BatchAdd(ctx context.Context, links []model.Link) ([]model.Link, error) {
	valueStrings := make([]string, 0, len(links))
	valueArgs := make([]any, 0, len(links)*3)
	for i, link := range links {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, link.UUID)
		valueArgs = append(valueArgs, link.ShortURL)
		valueArgs = append(valueArgs, link.OriginalURL)
		valueArgs = append(valueArgs, link.UserID)
	}
	query := fmt.Sprintf(`
		INSERT INTO links
		    (uuid, short_url, original_url, user_id)
		VALUES %s`,
		strings.Join(valueStrings, ","))
	res, err := s.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, err
	}

	if rows, _ := res.RowsAffected(); rows != int64(len(links)) {
		return nil, fmt.Errorf("insert error")
	}

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

func (s *dbStorage) GetByUser(ctx context.Context, userID int) ([]model.Link, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT short_url
		     , original_url
		FROM links
		WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	var links []model.Link
	for rows.Next() {
		var link model.Link
		if err := rows.Scan(&link.ShortURL, &link.OriginalURL); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func NewDBStorage(db *sql.DB) *dbStorage {
	return &dbStorage{
		db: db,
	}
}
