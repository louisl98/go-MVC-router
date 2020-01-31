package model

import (
	"database/sql"
	"errors"
)

var (
	// ErrCode is a config or an internal error
	ErrCode = errors.New("case statement in code is not correct")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("result not found")
	// ErrUnavailable is a database not available error
	ErrUnavailable = errors.New("database is unavailable")
	// ErrUnauthorized is a permissions violation
	ErrUnauthorized = errors.New("user does not have permission to perform this operation")
)

// StandardizeError returns the same error regardless of the database used
func StandardizeError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNoResult
	}
	return err
}

func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}
