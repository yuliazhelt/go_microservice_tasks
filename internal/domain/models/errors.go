package models

import "fmt"

var (
	ErrForbidden = fmt.Errorf("forbidden")
	ErrTokenExpired = fmt.Errorf("token expired")
	ErrTokenInvalid = fmt.Errorf("token invalid")
	ErrNotFound = fmt.Errorf("not found")
)