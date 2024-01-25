package web

import "errors"

var (
	ErrorVersionFileNotFound  = errors.New("release version file not found")
	ErrorVersionFileInvalid   = errors.New("release version file invalid")
	ErrorDatabaseFileNotFound = errors.New("database file not found")
	ErrorDatabaseFileInvalid  = errors.New("database file invalid")
)
