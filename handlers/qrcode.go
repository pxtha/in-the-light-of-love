package handlers

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

func QRCodeHandler(w http.ResponseWriter, r *http.Request) {
	url := "http://157.245.59.240:8080/"
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}
