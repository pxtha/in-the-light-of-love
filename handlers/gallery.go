package handlers

import (
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type GuestStats struct {
	Name  string
	Likes int
}

type PageData struct {
	Username    string
	Gallery     map[string][]Photo
	TotalGuests int
	TotalPhotos int
	TotalLikes  int
	BestPhoto   Photo
	TopGuests   []GuestStats
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Gallery(db *gorm.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("templates/gallery.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var photos []Photo
		db.Find(&photos)

		gallery := make(map[string][]Photo)
		for _, p := range photos {
			gallery[p.Uploader] = append(gallery[p.Uploader], p)
		}

		for _, p := range gallery {
			sort.Slice(p, func(i, j int) bool {
				if p[i].Likes != p[j].Likes {
					return p[i].Likes > p[j].Likes
				}
				return p[i].ModTime.After(p[j].ModTime)
			})
		}

		totalGuests := len(gallery)
		totalPhotos := 0
		totalLikes := 0
		for _, p := range gallery {
			totalPhotos += len(p)
			for _, photo := range p {
				totalLikes += photo.Likes
			}
		}

		// Find best photo and top guests of the day
		var bestPhoto Photo
		guestLikesToday := make(map[string]int)
		today := time.Now().Truncate(24 * time.Hour)

		for uploader, p := range gallery {
			for _, photo := range p {
				if photo.ModTime.After(today) {
					if photo.Likes > bestPhoto.Likes {
						bestPhoto = photo
					}
					guestLikesToday[uploader] += photo.Likes
				}
			}
		}

		var topGuests []GuestStats
		for name, likes := range guestLikesToday {
			topGuests = append(topGuests, GuestStats{Name: name, Likes: likes})
		}
		sort.Slice(topGuests, func(i, j int) bool {
			return topGuests[i].Likes > topGuests[j].Likes
		})
		if len(topGuests) > 2 {
			topGuests = topGuests[:2]
		}

		data := PageData{
			Username:    username,
			Gallery:     gallery,
			TotalGuests: totalGuests,
			TotalPhotos: totalPhotos,
			TotalLikes:  totalLikes,
			BestPhoto:   bestPhoto,
			TopGuests:   topGuests,
		}

		tmpl.Execute(w, data)
	}
}
