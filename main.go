package main

import (
	"fmt"
	"in-the-light-of-love/handlers"
	"log"
	"net/http"
	"os"
	"strings"
	_ "time"

	"github.com/gorilla/sessions"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("/app/gallery.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Drop the table to clear all data on restart, then recreate it.
	db.Migrator().DropTable(&handlers.Photo{})
	err = db.AutoMigrate(&handlers.Photo{})
	if err != nil {
		log.Fatal("failed to migrate database schema")
	}

	// Repopulate database from uploads folder
	repopulateDbFromUploads(db)

	store := sessions.NewCookieStore([]byte("super-secret-key"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/", handlers.LoginPageHandler)
	http.HandleFunc("/gallery", handlers.Gallery(db, store))
	http.HandleFunc("/upload", handlers.Upload(db, store))
	http.HandleFunc("/qr", handlers.QRCodeHandler)
	http.HandleFunc("/clear", handlers.ClearHandler)
	http.HandleFunc("/login", handlers.Login(store))
	http.HandleFunc("/logout", handlers.Logout(store))
	http.HandleFunc("/like", handlers.Like(db))
	http.HandleFunc("/delete", handlers.DeletePhoto(db, store))

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func repopulateDbFromUploads(db *gorm.DB) {
	uploadDir := "uploads"
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Printf("Could not read uploads directory: %v. This might be expected if the directory doesn't exist yet.", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			parts := strings.SplitN(fileName, "_", 2)
			if len(parts) == 2 {
				uploader := parts[0]

				fileInfo, err := file.Info()
				if err != nil {
					log.Printf("Could not get file info for %s: %v", fileName, err)
					continue
				}

				photo := handlers.Photo{
					Filename: fileName,
					Uploader: uploader,
					Likes:    0, // Likes are not stored in filenames, so they reset
					ModTime:  fileInfo.ModTime(),
				}
				db.Create(&photo)
			}
		}
	}
	log.Println("Successfully repopulated database from uploads folder.")
}
