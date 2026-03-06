package errs

import "errors"

var ErrTeamNotFound = errors.New("team not found")

var ErrTeamAlreadyExists = errors.New("team already exists")

var ErrTeamMemberNotFound = errors.New("team member not found")
var ErrUserAlreadyInTeam = errors.New("user already in team")

var ErrNoRightsToDeleteMember = errors.New("no rights to delete member")
var ErrNoRightsToInviteMember = errors.New("no rights to invite member")
