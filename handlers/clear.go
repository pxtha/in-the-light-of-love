package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func ClearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	password := r.FormValue("password")
	if password != "admm@123" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uploadDir := "uploads"
	dir, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, "Failed to read uploads directory", http.StatusInternalServerError)
		return
	}

	for _, d := range dir {
		os.RemoveAll(filepath.Join(uploadDir, d.Name()))
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
