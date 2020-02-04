package model

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

// generate random number
func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// FormUploadsGET gets all uploaded files in the form
func (p *Post) FormUploadsGET(w http.ResponseWriter, r *http.Request) {
	if f, _ := r.MultipartForm.File["upload"]; f != nil {
		for _, files := range r.MultipartForm.File {
			for i := range files {
				r.ParseMultipartForm(32 << 20)
				file, err := files[i].Open()
				defer file.Close()
				if err != nil {
					log.Println(w, err)
				}
				tempFile, filename, e := TempFile("uploads", files[i].Filename)
				if e != nil {
					log.Println(e)
				}
				defer tempFile.Close()
				fileBytes, ee := ioutil.ReadAll(file)
				if ee != nil {
					log.Println(ee)
				}
				tempFile.Write(fileBytes)
				filename = strings.Replace(filename, "uploads/", "", 1)
				UploadCreate(filename, p.ID)
			}
		}
	}
}

// TempFile uses a random number as prefix for file name to avoid file overwriting and returns a new file and its file name
func TempFile(dir, pattern string) (f *os.File, name string, err error) {
	if dir == "" {
		dir = os.TempDir()
	}
	var prefix, suffix string
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		suffix, prefix = pattern[:pos], pattern[pos+1:]
	} else {
		suffix = pattern
	}
	nconflict := 0
	for i := 0; i < 10000; i++ {
		name = filepath.Join(dir, prefix+nextRandom()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

// UploadCreate creates an upload in the database
func UploadCreate(filename string, postID uint32) error {
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
func FileDelete(ID string, postID string) (error, error) {
	r := Upload{}
	e := database.SQL.Get(&r, "SELECT id, file_name, short_name, post_id, uploaded_at, deleted FROM uploads WHERE id = ? AND post_id = ? LIMIT 1", ID, postID)
	filepath := r.FileName
	var ee = os.Remove("uploads/" + filepath)
	if ee != nil {
		log.Println(ee)
	}
	_, err := database.SQL.Exec("DELETE FROM UPLOADS WHERE id = ? AND post_id = ?", ID, postID)
	return StandardizeError(e), StandardizeError(err)
}
