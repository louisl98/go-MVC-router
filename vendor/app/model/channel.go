package model

import (
	"fmt"
	"log"
	"net/http"

	"time"

	"app/shared/database"
	"app/shared/session"
	"app/shared/view"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Channel
// *****************************************************************************

// Channel table contains the information for each channel
type Channel struct {
	ID        uint32        `db:"id" bson:"id,omitempty"` // Don't use Id, use UserID() instead for consistency with MongoDB
	Username  bson.ObjectId `db:"username" bson:"_username"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
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
			c := session.DB(database.ReadConfig().MongoDB.Database).C("user")

			// Validate the object username
			if bson.IsObjectIdHex(username) {
				err = c.FindId(bson.ObjectIdHex(username)).One(&result)
				if result.Username != bson.ObjectIdHex(username) {
					result = Channel{}
					err = ErrUnauthorized
				}
			} else {
				err = ErrNoResult
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		err = database.View("user", username, &result)
		if err != nil {
			err = ErrNoResult
		}
		if result.Username != bson.ObjectIdHex(username) {
			result = Channel{}
			err = ErrUnauthorized
		}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// ChannelReadGET displays the data on the channel page
func ChannelReadGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	username := fmt.Sprintf("%s", sess.Values["username"])

	channel, err := ChannelByUsername(username)
	if err != nil {
		log.Println(err)
		channel = Channel{}
	}

	// Display the view
	v := view.New(r)
	v.Name = "channel/channel"
	v.Vars["channel"] = channel
	v.Vars["username"] = username
	v.Vars["posts"] = posts
	v.Render(w)
}
