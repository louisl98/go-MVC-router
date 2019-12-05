package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"app/shared/database"

	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Post
// *****************************************************************************

// Post table contains the information for each post
type Post struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID        uint32        `db:"id" bson:"id,omitempty"` // Don't use Id, use PostID() instead for consistency with MongoDB
	Content   string        `db:"content" bson:"content"`
	UserID    bson.ObjectId `bson:"user_id"`
	UID       uint32        `db:"user_id" bson:"userid,omitempty"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
}

// PostID returns the post id
func (u *Post) PostID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		r = fmt.Sprintf("%v", u.ID)
	case database.TypeMongoDB:
		r = u.ObjectID.Hex()
	case database.TypeBolt:
		r = u.ObjectID.Hex()
	}

	return r
}

// PostByID gets post by ID
func PostByID(userID string, postID string) (Post, error) {
	var err error

	result := Post{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE id = ? AND user_id = ? LIMIT 1", postID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("post")

			// Validate the object id
			if bson.IsObjectIdHex(postID) {
				err = c.FindId(bson.ObjectIdHex(postID)).One(&result)
				if result.UserID != bson.ObjectIdHex(userID) {
					result = Post{}
					err = ErrUnauthorized
				}
			} else {
				err = ErrNoResult
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		err = database.View("post", userID+postID, &result)
		if err != nil {
			err = ErrNoResult
		}
		if result.UserID != bson.ObjectIdHex(userID) {
			result = Post{}
			err = ErrUnauthorized
		}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// PostsByUserID gets all posts for a user
func PostsByUserID(userID string) ([]Post, error) {
	var err error

	var result []Post

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM post WHERE user_id = ?", userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("post")

			// Validate the object id
			if bson.IsObjectIdHex(userID) {
				err = c.Find(bson.M{"user_id": bson.ObjectIdHex(userID)}).All(&result)
			} else {
				err = ErrNoResult
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		// View retrieves a record set in Bolt
		err = database.BoltDB.View(func(tx *bolt.Tx) error {
			// Get the bucket
			b := tx.Bucket([]byte("post"))
			if b == nil {
				return bolt.ErrBucketNotFound
			}

			// Get the iterator
			c := b.Cursor()

			prefix := []byte(userID)
			for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
				var single Post

				// Decode the record
				err := json.Unmarshal(v, &single)
				if err != nil {
					log.Println(err)
					continue
				}

				result = append(result, single)
			}

			return nil
		})
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// PostCreate creates a post
func PostCreate(content string, userID string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("INSERT INTO post (content, user_id) VALUES (?,?)", content, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("post")

			post := &Post{
				ObjectID:  bson.NewObjectId(),
				Content:   content,
				UserID:    bson.ObjectIdHex(userID),
				CreatedAt: now,
				UpdatedAt: now,
				Deleted:   0,
			}
			err = c.Insert(post)
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		post := &Post{
			ObjectID:  bson.NewObjectId(),
			Content:   content,
			UserID:    bson.ObjectIdHex(userID),
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}

		err = database.Update("post", userID+post.ObjectID.Hex(), &post)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// PostUpdate updates a post
func PostUpdate(content string, userID string, postID string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("UPDATE post SET content=? WHERE id = ? AND user_id = ? LIMIT 1", content, postID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("post")
			var post Post
			post, err = PostByID(userID, postID)
			if err == nil {
				// Confirm the owner is attempting to modify the post
				if post.UserID.Hex() == userID {
					post.UpdatedAt = now
					post.Content = content
					err = c.UpdateId(bson.ObjectIdHex(postID), &post)
				} else {
					err = ErrUnauthorized
				}
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		var post Post
		post, err = PostByID(userID, postID)
		if err == nil {
			// Confirm the owner is attempting to modify the post
			if post.UserID.Hex() == userID {
				post.UpdatedAt = now
				post.Content = content
				err = database.Update("post", userID+post.ObjectID.Hex(), &post)
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// PostDelete deletes a post
func PostDelete(userID string, postID string) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("DELETE FROM post WHERE id = ? AND user_id = ?", postID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("post")

			var post Post
			post, err = PostByID(userID, postID)
			if err == nil {
				// Confirm the owner is attempting to modify the post
				if post.UserID.Hex() == userID {
					err = c.RemoveId(bson.ObjectIdHex(postID))
				} else {
					err = ErrUnauthorized
				}
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		var post Post
		post, err = PostByID(userID, postID)
		if err == nil {
			// Confirm the owner is attempting to modify the post
			if post.UserID.Hex() == userID {
				err = database.Delete("post", userID+post.ObjectID.Hex())
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
