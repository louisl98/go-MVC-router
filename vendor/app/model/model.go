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

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNoResult
	}
	return err
}
