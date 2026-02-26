package validation

import (
	"errors"
	"net/url"
)

var ErrURINotAbsolute = errors.New("uri is not absolute")

func IsURI(value string) (*url.URL, error) {
	uri, err := url.Parse(value)
	if err != nil {
		return nil, err
	}
	if !uri.IsAbs() {
		return nil, ErrURINotAbsolute
	}
	return uri, nil
}
