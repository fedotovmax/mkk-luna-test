package errs

import "errors"

var ErrUserNotFound = errors.New("user not found")

var ErrBadCredentials = errors.New("bad credentials")

var ErrUserAlreadyExists = errors.New("user already exists")
