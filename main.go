package main

import (
	"fmt"
	"in-the-light-of-love/handlers"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("gallery.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&handlers.Photo{})
	if err != nil {
		log.Fatal("failed to migrate database schema")
	}

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
