package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PageData struct {
	Gallery map[string][]string
}

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gallery := make(map[string][]string)
	uploadDir := "uploads"

	err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			parts := strings.Split(info.Name(), "_")
			if len(parts) > 1 {
				uploaderName := parts[0]
				gallery[uploaderName] = append(gallery[uploaderName], info.Name())
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Gallery: gallery,
	}

	tmpl.Execute(w, data)
}
