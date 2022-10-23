package internal

import "errors"

var (
	ErrRecognizedMigrationType = errors.New("not a recognized migration file type")
	ErrNoFilenameSeparator     = errors.New("no filename separator '_' found")
	ErrInvalidMigrationId      = errors.New("migration IDs must be greater than zero")
)
