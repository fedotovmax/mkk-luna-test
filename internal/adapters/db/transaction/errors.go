package transaction

import "errors"

var ErrConnRequired = errors.New("connection is required for create transaction manager")
