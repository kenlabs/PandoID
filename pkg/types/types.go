package types

import (
	"encoding/json"
	"net/url"
)

type URI struct {
	url.URL
}

func (u URI) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u URI) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u URI) UnmarshalJSON(bytes []byte) error {
	var value string

	err := json.Unmarshal(bytes, &value)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return err
	}
	u.URL = *parsedURL

	return nil
}

func ParseURI(uri string) (*URI, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return &URI{*parsedURI}, nil
}

func (u URI) String() string {
	return u.URL.String()
}
