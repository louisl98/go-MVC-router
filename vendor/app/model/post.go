package model

import (
	"time"

	"app/shared/database"
)

// *****************************************************************************
// Post
// *****************************************************************************

// Post struct contains the information for each post
type Post struct {
	ID        uint32    `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UserID    uint32    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   uint8     `db:"deleted"`
	Files     []Upload
}

// PostByID gets a post by ID
func PostByID(postID string, userID string) (Post, error) {
	result := Post{}
	err := database.SQL.Get(&result, "SELECT id, title, content, user_id, created_at, updated_at, deleted FROM post WHERE id = ? AND user_id = ? LIMIT 1", postID, userID)
	result.UploadsGET()
	return result, StandardizeError(err)
}

// PostsByUserID gets all posts for a user
func PostsByUserID(userID string) ([]Post, error) {
	var result []Post
	err := database.SQL.Select(&result, "SELECT id, title, content, user_id, created_at, updated_at, deleted FROM post WHERE user_id = ?", userID)
	for r := range result {
		result[r].UploadsGET()
	}
	return result, StandardizeError(err)
}

// PostCreate creates a post and returns it
func PostCreate(title string, content string, userID string) (Post, error, error) {
	_, e := database.SQL.Exec("INSERT INTO post (title, content, user_id) VALUES (?,?,?)", title, content, userID)
	result := Post{}
	err := database.SQL.Get(&result, "SELECT id, created_at, updated_at, deleted FROM post WHERE title = ? AND content = ? AND user_id = ? LIMIT 1", title, content, userID)
	return result, StandardizeError(e), StandardizeError(err)
}

// PostUpdate updates a post and returns it
func PostUpdate(title string, content string, userID string, postID string) (Post, error, error) {
	_, e := database.SQL.Exec("UPDATE post SET title=?, content=? WHERE user_id = ? AND id = ? LIMIT 1", title, content, userID, postID)
	result := Post{}
	err := database.SQL.Get(&result, "SELECT id, title, content, user_id, created_at, updated_at, deleted FROM post WHERE id = ? AND user_id = ? LIMIT 1", postID, userID)
	return result, StandardizeError(e), StandardizeError(err)
}

// PostDelete deletes a post
func PostDelete(postID string, userID string) error {
	_, err := database.SQL.Exec("DELETE FROM post WHERE id = ? AND user_id = ?", postID, userID)
	return StandardizeError(err)
}
