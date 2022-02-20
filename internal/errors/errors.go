package errors

import "fmt"

var (
	InvalidEmailOrPassword = fmt.Errorf("invalid email or password")
)
