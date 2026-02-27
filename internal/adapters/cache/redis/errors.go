package redis

import "errors"

var ErrCloseTimeout = errors.New("the time to safely terminate the connection to the redis has expired")
