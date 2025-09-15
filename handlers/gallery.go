package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PageData struct {
	Username    string
	Gallery     map[string][]Photo
	TotalGuests int
	TotalPhotos int
	TotalLikes  int
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	username := ""
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username = cookie.Value

	tmpl, err := template.ParseFiles("templates/gallery.html")
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
					ModTime:  info.ModTime(),
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

	for _, photos := range gallery {
		sort.Slice(photos, func(i, j int) bool {
			if photos[i].Likes != photos[j].Likes {
				return photos[i].Likes > photos[j].Likes
			}
			return photos[i].ModTime.After(photos[j].ModTime)
		})
	}

	totalGuests := len(gallery)
	totalPhotos := 0
	totalLikes := 0
	for _, photos := range gallery {
		totalPhotos += len(photos)
		for _, photo := range photos {
			totalLikes += photo.Likes
		}
	}

	data := PageData{
		Username:    username,
		Gallery:     gallery,
		TotalGuests: totalGuests,
		TotalPhotos: totalPhotos,
		TotalLikes:  totalLikes,
	}

	tmpl.Execute(w, data)
}
