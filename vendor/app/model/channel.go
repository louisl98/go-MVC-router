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

// ChannelByUsername gets channel information from username
func ChannelByUsername(username string) (Channel, error) {
	var err error
	result := Channel{}
	err = database.SQL.Get(&result, "SELECT id, username, email, created_at, updated_at FROM user WHERE username = ? LIMIT 1", username)
	return result, standardizeError(err)
}

// PostsByChannelID gets all posts for a channel
func PostsByChannelID(channelID uint32) ([]Post, error) {
	var err error
	var result []Post
	err = database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE user_id = ?", channelID)
	return result, standardizeError(err)
}

// ChannelReadGET gets the query and displays the channel
func ChannelReadGET(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	var request = strings.Trim(uri, "/channel/")
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
<<<<<<< HEAD

// ChannelByUsername gets the channel by username
func ChannelByUsername(username string) (Channel, error) {
	var err error

	result := Channel{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, username, created_at FROM user WHERE username = ?  LIMIT 1", username)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			result = Channel{}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		err = database.View("user", username, &result)
		if err != nil {
			err = ErrNoResult
		}
		result = Channel{}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}
=======
>>>>>>> 9c7a11885182d87ab8517d8575d2634674825124
