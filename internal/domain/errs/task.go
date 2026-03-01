package errs

import "errors"

var ErrTaskNotFound = errors.New("task not found")

var ErrNoRightsToDeleteTaskComment = errors.New("no rights to delete task comment")

var ErrUserNotInTaskTeam = errors.New("user is not in prepared task's team")
