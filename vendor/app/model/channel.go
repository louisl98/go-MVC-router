package model

import (
	"app/shared/database"
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
	CreatedAt time.Time `db:"created_at" bson:"created_at"`
}

// ChannelReadGET displays the data on the channel page
func ChannelReadGET(w http.ResponseWriter, r *http.Request) {

	request := r.RequestURI
	var username = strings.Trim(request, "/channel/")

	channel, err := ChannelByUsername(username)

	if err != nil {
		log.Println(err)
		return
	}

	channel = Channel{}

	// Display the view
	v := view.New(r)
	v.Name = "channel/channel"
	v.Vars["channel"] = channel
	v.Vars["username"] = username
	v.Render(w)
}

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
