package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// where user has invalid auth
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// where user tries to signup with used email
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
