package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ClearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	name := cookie.Value

	uploadDir := "uploads"
	dir, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, "Failed to read uploads directory", http.StatusInternalServerError)
		return
	}

	for _, d := range dir {
		if strings.HasPrefix(d.Name(), name+"_") {
			os.RemoveAll(filepath.Join(uploadDir, d.Name()))
		}
	}

	http.Redirect(w, r, "/gallery", http.StatusSeeOther)
}
