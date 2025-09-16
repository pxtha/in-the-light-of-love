package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func Login(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		session, _ := store.Get(r, "session-name")
		session.Values["username"] = name
		session.Save(r, w)

		http.Redirect(w, r, "/gallery", http.StatusSeeOther)
	}
}

func Logout(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		session.Values["username"] = ""
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
