package model

import (
	"time"

	"app/shared/database"
)

// Upload struct contains the information for each post
type Upload struct {
	ID         uint32    `db:"id"`
	FileName   string    `db:"file_name"`
	PostID     uint32    `db:"post_id"`
	UploadedAt time.Time `db:"uploaded_at"`
	Deleted    uint8     `db:"deleted"`
}

// UploadCreate creates an upload in the database
func UploadCreate(filename string, postID string) error {
	_, err := database.SQL.Exec("INSERT INTO uploads (file_name, post_id) VALUES (?,?)", filename, postID)
	return StandardizeError(err)
}

// UploadsByPostID gets all uploads for a post
func UploadsByPostID(postID uint32) ([]Upload, error) {
	var result []Upload
	err := database.SQL.Select(&result, "SELECT id, file_name, post_id, uploaded_at, deleted FROM uploads WHERE post_id = ?", postID)
	return result, StandardizeError(err)
}

// UploadDelete deletes an upload
func UploadDelete(ID string, postID string) error {
	_, err := database.SQL.Exec("DELETE FROM uploads WHERE id = ? AND post_id = ?", ID, postID)
	return StandardizeError(err)
}