package errors

import "fmt"

var (
	ErrInvalidEmailOrPassword = fmt.Errorf("invalid email or password")
	ErrUserNotInGroup         = fmt.Errorf("user is not in a group")
	ErrUserNotMemberOfGroup   = fmt.Errorf("user is not a member of the group")
)
