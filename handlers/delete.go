package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func DeletePhoto(db *gorm.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		filename := r.FormValue("filename")
		if filename == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		var photo Photo
		if err := db.Where("filename = ?", filename).First(&photo).Error; err != nil {
			http.Error(w, "Photo not found", http.StatusNotFound)
			return
		}

		if photo.Uploader != username {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Delete the photo file
		filePath := filepath.Join("uploads", filename)
		if err := os.Remove(filePath); err != nil {
			http.Error(w, "Failed to delete photo file", http.StatusInternalServerError)
			return
		}

		// Delete the photo record from the database
		if err := db.Delete(&photo).Error; err != nil {
			http.Error(w, "Failed to delete photo record", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
