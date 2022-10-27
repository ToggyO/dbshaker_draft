package internal

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

func IsValidFileName(value string) (int64, error) {
	base := filepath.Base(value)
	if ext := filepath.Ext(base); ext != GoExt && ext != SqlExt {
		return 0, ErrRecognizedMigrationType
	}

	index := strings.Index(base, FileNameSeparator)
	if index < 0 {
		return 0, ErrNoFilenameSeparator
	}

	// TODO:  разобраться, как работает ParseInt
	num, err := strconv.ParseInt(base[:index], 10, 64)
	if err == nil && num <= 0 {
		return 0, ErrInvalidMigrationId
	}

	return num, nil
}

func LogWithPrefix(message string) {
	log.Printf("%s: %s", ToolName, message)
}
