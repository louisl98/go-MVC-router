package model

import (
	"time"

	"app/shared/database"
)

// *****************************************************************************
// Post
// *****************************************************************************

// Post table contains the information for each post
type Post struct {
	ID        uint32    `db:"id"`
	Content   string    `db:"content"`
	UserID    uint32    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   uint8     `db:"deleted"`
}

// PostByID gets post by ID
func PostByID(postID string, userID string) (Post, error) {
	var err error
	result := Post{}
	err = database.SQL.Get(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE id = ? AND user_id = ? LIMIT 1", postID, userID)
	return result, StandardizeError(err)
}

// PostsByUserID gets all posts for a user
func PostsByUserID(userID string) ([]Post, error) {
	var err error
	var result []Post
	err = database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE user_id = ?", userID)
	return result, StandardizeError(err)
}

// PostCreate creates a post
func PostCreate(content string, userID string) error {
	var err error
	_, err = database.SQL.Exec("INSERT INTO post (content, user_id) VALUES (?,?)", content, userID)
	return StandardizeError(err)
}

// PostUpdate updates a post
func PostUpdate(content string, userID string, postID string) error {
	var err error
	_, err = database.SQL.Exec("UPDATE post SET content=? WHERE user_id = ? AND id = ? LIMIT 1", content, userID, postID)
	return StandardizeError(err)
}

// PostDelete deletes a post
func PostDelete(postID string, userID string) error {
	var err error
	_, err = database.SQL.Exec("DELETE FROM post WHERE id = ? AND user_id = ?", postID, userID)
	return StandardizeError(err)
}

func GetNextPostID() (uint32, error){
	var err error
	result := Post{}
	err = database.SQL.Get(&result, "SELECT id FROM post ORDER BY ID DESC LIMIT 1")
	return result.ID, StandardizeError(err)
}