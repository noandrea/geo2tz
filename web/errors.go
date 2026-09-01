package web

import "errors"

var (
	ErrVersionFileNotFound  = errors.New("release version file not found")
	ErrVersionFileInvalid   = errors.New("release version file invalid")
	ErrDatabaseFileNotFound = errors.New("database file not found")
	ErrDatabaseFileInvalid  = errors.New("database file invalid")
)
