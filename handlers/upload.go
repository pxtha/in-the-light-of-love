package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func Upload(db *gorm.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseMultipartForm(10 << 20) // 10 MB

		session, _ := store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		files := r.MultipartForm.File["photos"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			fileName := fmt.Sprintf("%s_%s", username, fileHeader.Filename)
			filePath := filepath.Join("uploads", fileName)

			dst, err := os.Create(filePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			photo := Photo{
				Filename: fileName,
				Uploader: username,
				Likes:    0,
				ModTime:  time.Now(),
			}
			db.Create(&photo)
		}

		http.Redirect(w, r, "/gallery", http.StatusSeeOther)
	}
}
