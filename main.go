package main

import (
	"fmt"
	"in-the-light-of-love/handlers"
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/", handlers.LoginPageHandler)
	http.HandleFunc("/gallery", handlers.GalleryHandler)
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/qr", handlers.QRCodeHandler)
	http.HandleFunc("/clear", handlers.ClearHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/like", handlers.LikeHandler)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
