package controller

import (
	"fmt"
	"log"
	"net/http"

	"app/model"
	"app/shared/session"
	"app/shared/view"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// ChannelReadGET displays the posts in the channel
func ChannelReadGET(w http.ResponseWriter, r *http.Request) {
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
	v.Name = "channel/channelcontent"
	v.Vars["username"] = sess.Values["username"]
	v.Vars["posts"] = posts
	v.Render(w)
}

// ChannelCreateGET displays the post creation page
func ChannelCreateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "channel/newpost"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

// ChannelCreatePOST handles the post creation form submission
func ChannelCreatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"post"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		ChannelCreateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("post")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get database result
	err := model.PostCreate(content, userID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Post added!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/channel", http.StatusFound)
		return
	}

	// Display the same page
	ChannelCreateGET(w, r)
}

// ChannelUpdateGET displays the post update page
func ChannelUpdateGET(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/channel", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "channel/editpost"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["post"] = post.Content
	v.Render(w)
}

// ChannelUpdatePOST handles the post update form submission
func ChannelUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"post"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		ChannelUpdateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("post")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	postID := params.ByName("id")

	// Get database result
	err := model.PostUpdate(content, userID, postID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Post updated!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/channel", http.StatusFound)
		return
	}

	// Display the same page
	ChannelUpdateGET(w, r)
}

// ChannelDeleteGET handles the post deletion
func ChannelDeleteGET(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "/channel", http.StatusFound)
	return
}
