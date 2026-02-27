package config

import "errors"

var ErrInvalidAppEnv = errors.New("app env is invalid or not supported")
var ErrConfigPathNotExists = errors.New("config path for dev env is not exists")
