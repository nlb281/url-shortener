package storage

import (
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
) 

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists = errors.New("url exists")
	ErrAliasExists = errors.New("alias exists")
	ErrURLNotDeleted = errors.New("failed to delete url")
)

func IsUniqueConstraint(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
	}
	return false
}