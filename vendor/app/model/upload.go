package model

import (
	"log"
	"time"

	"app/shared/database"
)

// Upload struct contains the information for each post
type Upload struct {
	ID         uint32    `db:"id"`
	FileName   string    `db:"file_name"`
	ShortName  string    `db:"short_name"`
	PostID     uint32    `db:"post_id"`
	UploadedAt time.Time `db:"uploaded_at"`
	Deleted    uint8     `db:"deleted"`
}

// UploadCreate creates an upload in the database
func UploadCreate(filename string, postID string) error {
	shortname := trimLeftChars(filename, 9)
	_, err := database.SQL.Exec("INSERT INTO uploads (file_name, short_name, post_id) VALUES (?,?,?)", filename, shortname, postID)
	return StandardizeError(err)
}

// UploadsGET gets all uploads for a post
func (p *Post) UploadsGET() {
	var result []Upload
	err := database.SQL.Select(&result, "SELECT id, file_name, short_name, post_id, uploaded_at, deleted FROM uploads WHERE post_id = ?", p.ID)
	if err != nil {
		log.Println(StandardizeError(err))
	}
	p.Files = result
}

// PostIDByFileID gets post ID by file ID
func PostIDByFileID(fileID string) (string, error) {
	var result string
	err := database.SQL.Get(&result, "SELECT post_id FROM uploads WHERE id = ? LIMIT 1", fileID)
	return result, StandardizeError(err)
}

// FileDelete deletes a file
func FileDelete(ID string, postID string) error {
	_, err := database.SQL.Exec("DELETE FROM UPLOADS WHERE id = ? AND post_id = ?", ID, postID)
	return StandardizeError(err)
}
