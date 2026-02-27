package validation

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidUUIDFormat = errors.New("invalid uuid format")

var ErrURINotAbsolute = errors.New("uri is not absolute")

var (
	ErrEmptyPath   = errors.New("path is empty")
	ErrInvalidPath = errors.New("invalid file path")
)

func IsUUID(value string) (uuid.UUID, error) {

	uid, err := uuid.Parse(value)

	if err != nil {
		return uuid.Nil, ErrInvalidUUIDFormat
	}

	return uid, nil
}

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

func IsFilePath(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return ErrEmptyPath
	}

	if strings.ContainsAny(value, "<>:\"|?*") {
		return ErrInvalidPath
	}

	if filepath.Base(value) == "." || filepath.Base(value) == ".." {
		return ErrInvalidPath
	}

	return nil
}
