// Package model
package model

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
)

type URL struct{ url.URL }

func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URL) UnmarshalJSON(data []byte) error {
	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	v, err := url.Parse(tmp)
	if err == nil {
		*u = URL{*v}
	}
	return err
}

func (u URL) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *URL) Scan(value any) error {
	if value == nil {
		*u = URL{}
		return nil
	}

	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return errDBScan
	}

	tmp, ok := sv.(string)
	if !ok {
		return errDBConvert
	}

	v, err := url.Parse(tmp)
	if err == nil {
		*u = URL{*v}
	}
	return err
}

type Link struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL URL    `json:"original_url"`
	UserID      int
	DeletedFlag bool
}

func NewLink(original url.URL, short string, userID int) Link {
	return Link{
		UUID:        short,
		ShortURL:    short,
		OriginalURL: URL{original},
		UserID:      userID,
		DeletedFlag: false,
	}
}

func NewLinkWithUUID(original url.URL, short string, userID int, uuid string) Link {
	return Link{
		UUID:        uuid,
		ShortURL:    short,
		OriginalURL: URL{original},
		UserID:      userID,
		DeletedFlag: false,
	}
}
