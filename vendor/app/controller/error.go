package controller

import (
	"app/shared/session"
	"app/shared/view"
	"fmt"
	"net/http"
)

// Error404 handles 404 - Page Not Found
func Error404(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "404/404"
	sess := session.Instance(r)
	v.Vars["username"] = sess.Values["username"]
	v.Render(w)
}

// Error500 handles 500 - Internal Server Error
func Error500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "Internal Server Error 500")
}

// InvalidToken handles CSRF attacks
func InvalidToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, `Your token <strong>expired</strong>, click <a href="javascript:void(0)" onclick="location.replace(document.referrer)">here</a> to try again.`)
}
