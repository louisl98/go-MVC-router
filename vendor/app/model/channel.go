package model

import (
	"app/shared/database"
	"app/shared/session"
	"app/shared/view"
	"log"
	"net/http"
	"strings"
	"time"
)

// *****************************************************************************
// Channel
// *****************************************************************************

// Channel table contains the information for each channel
type Channel struct {
	ID        uint32    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ChannelByUsername gets the channel information from username
func ChannelByUsername(username string) (Channel, error) {
	var err error
	result := Channel{}
	err = database.SQL.Get(&result, "SELECT id, username, created_at FROM user WHERE username = ?  LIMIT 1", username)
	return result, StandardizeError(err)
}

// PostsByChannelID gets all posts for a channel
func PostsByChannelID(channelID uint32) ([]Post, error) {
	var err error
	var result []Post
	err = database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE user_id = ?", channelID)
	for r := range result {
		result[r].UploadsGET()
	}
	return result, StandardizeError(err)
}

// ChannelReadGET gets the query and displays the channel
func ChannelReadGET(w http.ResponseWriter, r *http.Request) {
	request := strings.Trim(r.RequestURI, "/channel/")
	// Check if user exists
	channel, err := ChannelByUsername(request)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/404", 301)
	}
	posts, err := PostsByChannelID(channel.ID)
	if err != nil {
		log.Println(err)
	}
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "channel/channel"
	v.Vars["title"] = request
	v.Vars["username"] = sess.Values["username"]
	v.Vars["creationdate"] = channel.CreatedAt
	v.Vars["lastseen"] = channel.UpdatedAt
	v.Vars["posts"] = posts
	v.Render(w)
}
