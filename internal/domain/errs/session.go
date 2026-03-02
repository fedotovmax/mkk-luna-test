package errs

import "errors"

var ErrSessionExpired = errors.New("session expired")
var ErrSessionNotFound = errors.New("session not found")
