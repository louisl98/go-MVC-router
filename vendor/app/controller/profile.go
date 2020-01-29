package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"app/model"
	"app/shared/session"
	"app/shared/view"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// ProfileReadGET displays the posts in the profile
func ProfileReadGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	userID := fmt.Sprintf("%s", sess.Values["id"])
	posts, err := model.PostsByUserID(userID)
	if err != nil {
		log.Println(err)
		posts = []model.Post{}
	}
	// Display the view
	v := view.New(r)
	v.Name = "profile/manageposts"
	v.Vars["username"] = sess.Values["username"]
	v.Vars["posts"] = posts
	v.Render(w)
}

// ProfileCreateGET displays the post creation page
func ProfileCreateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "profile/newpost"
	v.Vars["username"] = sess.Values["username"]
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

// ProfileCreatePOST handles the post creation form submission
func ProfileCreatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"post"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		ProfileCreateGET(w, r)
		return
	}
	// Get form values
	content := r.FormValue("post")
	userID := fmt.Sprintf("%s", sess.Values["id"])
	// Get uploaded file
	r.ParseMultipartForm(10 << 20)
	file, handler, e := r.FormFile("upload")
	if e != nil {
		fmt.Println(e)
		return
	}
	// add random prefix to filename and upload it to folder
	defer file.Close()
	tempFile, filename, ee := TempFile("uploads", handler.Filename)
	if ee != nil {
		fmt.Println(ee)
	}
	defer tempFile.Close()
	fileBytes, eee := ioutil.ReadAll(file)
	if eee != nil {
		fmt.Println(eee)
	}
	tempFile.Write(fileBytes)
	filename = strings.Replace(filename, "uploads/", "", 1)
	err := model.PostCreate(content, filename, userID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Post added!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	// Display the same page
	ProfileCreateGET(w, r)
}

// ProfileUpdateGET displays the post update page
func ProfileUpdateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Get the post id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	postID := params.ByName("id")
	userID := fmt.Sprintf("%s", sess.Values["id"])
	// Get the post
	post, err := model.PostByID(userID, postID)
	if err != nil { // If the post doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	// Display the view
	v := view.New(r)
	v.Name = "profile/editpost"
	v.Vars["username"] = sess.Values["username"]
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["post"] = post.Content
	v.Vars["file"] = post.FileName
	v.Render(w)
}

// ProfileUpdatePOST handles the post update form submission
func ProfileUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"post"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		ProfileUpdateGET(w, r)
		return
	}
	// Get form values
	content := r.FormValue("post")
	userID := fmt.Sprintf("%s", sess.Values["id"])
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	postID := params.ByName("id")
	// Get uploaded file
	r.ParseMultipartForm(10 << 20)
	file, handler, e := r.FormFile("upload")
	if e != nil {
		fmt.Println(e)
		return
	}
	// add random prefix to filename and upload it to folder
	defer file.Close()
	tempFile, filename, ee := TempFile("uploads", handler.Filename)
	if ee != nil {
		fmt.Println(ee)
	}
	defer tempFile.Close()
	fileBytes, eee := ioutil.ReadAll(file)
	if eee != nil {
		fmt.Println(eee)
	}
	tempFile.Write(fileBytes)
	filename = strings.Replace(filename, "uploads/", "", 1)
	err := model.PostUpdate(content, filename, userID, postID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Post updated!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	// Display the same page
	ProfileUpdateGET(w, r)
}

// ProfileDeleteGET handles the post deletion
func ProfileDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	userID := fmt.Sprintf("%s", sess.Values["id"])
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	postID := params.ByName("id")
	// Get database result
	err := model.PostDelete(userID, postID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Post deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}
	http.Redirect(w, r, "/profile", http.StatusFound)
	return
}
