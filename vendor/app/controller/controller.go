package controller

import (
	"app/model"
	"app/shared/session"
	"app/shared/view"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

// TempFile uses random number as prefix for file name to avoid file overwriting
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

// UploadServe maps static files
func UploadServe(w http.ResponseWriter, r *http.Request) {
	// Disable listing directories
	if strings.HasSuffix(r.URL.Path, "/") {
		Error404(w, r)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}

// FileDeleteGET handles the file deletion
func FileDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	ID := strings.Trim(r.RequestURI, "/profile/deletefile/")
	log.Println(ID)
	postID, e := model.PostIDByFileID(ID)
	if e != nil {
		log.Println(e)
	}
	// Get database result
	err := model.FileDelete(ID, postID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"File deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}
	http.Redirect(w, r, "/profile/editpost/"+postID, http.StatusFound)
	return
}
