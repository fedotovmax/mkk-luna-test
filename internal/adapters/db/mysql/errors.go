package mysql

import "errors"

var ErrWantToCallMethodsAfterInitPool = errors.New("you will be able to call mysql pool methods only after the connection has been created and initialized")

var ErrInvalidDSNFormat = errors.New("invalid mysql dsn format")
