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
	ID        uint32    `db:"id" bson:"id,omitempty"` // Don't use Id, use UserID() instead for consistency with MongoDB
	Username  string    `db:"username" bson:"username"`
	Email     string    `db:"email" bson:"email"`
	CreatedAt time.Time `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time `db:"updated_at" bson:"updated_at"`
}

// ChannelByUsername gets channel information from username
func ChannelByUsername(username string) (Channel, error) {
	var err error
	result := Channel{}
	err = database.SQL.Get(&result, "SELECT id, username, email, created_at, updated_at FROM user WHERE username = ? LIMIT 1", username)
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
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "channel/channel"
	v.Vars["username"] = request
	v.Vars["loggeduser"] = sess.Values["username"]
	v.Vars["creationdate"] = channel.CreatedAt
	v.Vars["lastseen"] = channel.UpdatedAt
	v.Render(w)
}
