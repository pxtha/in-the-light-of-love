package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PageData struct {
	Username string
	Gallery  map[string][]Photo
}

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gallery := make(map[string][]Photo)
	uploadDir := "uploads"

	err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			parts := strings.Split(info.Name(), "_")
			if len(parts) > 1 {
				uploaderName := parts[0]
				likesMutex.Lock()
				likes := photoLikes[info.Name()]
				likesMutex.Unlock()
				photo := Photo{
					Filename: info.Name(),
					Likes:    likes,
				}
				gallery[uploaderName] = append(gallery[uploaderName], photo)
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := ""
	cookie, err := r.Cookie("username")
	if err == nil {
		username = cookie.Value
	}

	data := PageData{
		Username: username,
		Gallery:  gallery,
	}

	tmpl.Execute(w, data)
}
