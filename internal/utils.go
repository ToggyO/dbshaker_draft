package internal

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

//func IsValidMigrationName(value string) (int64, error) {
//	index := strings.Index(value, FileNameSeparator)
//	if index < 0 {
//		return 0, ErrNoFilenameSeparator
//	}
//
//	num, err := strconv.ParseInt(value[index:], 10, 64)
//	if err == nil && num <= 0 {
//		return 0, ErrInvalidMigrationId
//	}
//
//	return num, nil
//}

// TODO: remove
func IsValidFileName(value string) (int64, error) {
	base := filepath.Base(value)
	if ext := filepath.Ext(base); ext != GoExt && ext != SqlExt {
		return 0, ErrRecognizedMigrationType
	}

	index := strings.Index(base, FileNameSeparator)
	if index < 0 {
		return 0, ErrNoFilenameSeparator
	}

	num, err := strconv.ParseInt(base[:index], 10, 64)
	if err == nil && num <= 0 {
		return 0, ErrInvalidMigrationId
	}

	return num, nil
}

func GetSuccessMigrationMessage(currentDbVersion int64) string {
	return fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion)
}
