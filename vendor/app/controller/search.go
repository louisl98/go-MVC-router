package controller

import (
	"net/http"

	"app/shared/session"
	"app/shared/view"
)

// SearchGET displays the About page
func SearchGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "search/search"
	v.Vars["username"] = sess.Values["username"]
	v.Render(w)
}
