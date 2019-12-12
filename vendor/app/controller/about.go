package controller

import (
	"net/http"

	"app/shared/session"
	"app/shared/view"
)

// AboutGET displays the About page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
	v.Vars["username"] = sess.Values["username"]
	v.Render(w)
}
