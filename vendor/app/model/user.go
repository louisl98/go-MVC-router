package model

import (
	"fmt"
	"time"

	"app/shared/database"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ID        uint32    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	StatusID  uint8     `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   uint8     `db:"deleted"`
}

// UserStatus table contains every possible user status (active/inactive)
type UserStatus struct {
	ID        uint8
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   uint8
}

// UserID returns the user id
func (u *User) UserID() string {
	r := ""
	r = fmt.Sprintf("%v", u.ID)
	return r
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	var err error
	result := User{}
	err = database.SQL.Get(&result, "SELECT id, password, status_id, username FROM user WHERE email = ? LIMIT 1", email)
	return result, StandardizeError(err)
}

// UserCreate creates user
func UserCreate(username, email, password string) error {
	var err error
	_, err = database.SQL.Exec("INSERT INTO user (username, email, password) VALUES (?,?,?)", username, email, password)
	return StandardizeError(err)
}
