package errs

import "errors"

var ErrTeamNotFound = errors.New("team not found")

var ErrTeamMemberNotFound = errors.New("team member not found")

var ErrNoRightsToDeleteMember = errors.New("no rights to delete member")
var ErrNoRightsToInviteMember = errors.New("no rights to invite member")
