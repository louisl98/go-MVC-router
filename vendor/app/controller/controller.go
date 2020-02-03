package controller

import (
	"app/model"
	"app/shared/session"
	"app/shared/view"
	"log"
	"net/http"
	"strings"
)

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
