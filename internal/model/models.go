// Package model
package model

import (
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

type Link struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL URL    `json:"original_url"`
}

func NewLink(original url.URL, short string) Link {
	return Link{
		UUID:        short,
		ShortURL:    short,
		OriginalURL: URL{original},
	}
}
