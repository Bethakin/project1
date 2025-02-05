package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// GET endpoint'i tanımla
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Query parametresini al
		word := r.URL.Query().Get("word")

		// Eğer "word" parametresi yoksa hata döndür
		if word == "" {
			w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
			json.NewEncoder(w).Encode(map[string]string{"error": "Lütfen 'word' parametresini ekleyin."})
			return
		}

		// JSON objesi oluştur ve döndür
		response := map[string]string{"word": word}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Sunucuyu başlat
	fmt.Println("Sunucu http://localhost:8080 adresinde çalışıyor...")
	http.ListenAndServe(":8080", nil)
}
