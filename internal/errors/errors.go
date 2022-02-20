package errors

import "fmt"

var (
	ErrEmailOrPasswordIncorrect = fmt.Errorf("invalid email or password")
)
