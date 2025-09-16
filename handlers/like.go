package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func Like(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		filename := r.FormValue("filename")
		if filename == "" {
			http.Error(w, "Filename is required", http.StatusBadRequest)
			return
		}

		var photo Photo
		db.First(&photo, "filename = ?", filename)
		photo.Likes++
		db.Save(&photo)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"likes": photo.Likes})
	}
}
