package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Kendi API'nizin GET endpoint'i
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Query parametresini al
		word := r.URL.Query().Get("word")

		// Eğer "word" parametresi yoksa hata döndür
		if word == "" {
			w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
			json.NewEncoder(w).Encode(map[string]string{"error": "Lütfen 'word' parametresini ekleyin."})
			return
		}

		// Hedef URL'ye istek gönder
		targetURL := fmt.Sprintf("https://hedef-site.com/arama?q=%s", word)
		resp, err := http.Get(targetURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
			json.NewEncoder(w).Encode(map[string]string{"error": "Hedef siteye istek gönderilirken hata oluştu."})
			return
		}
		defer resp.Body.Close()

		// Hedef siteden gelen yanıtı oku
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
			json.NewEncoder(w).Encode(map[string]string{"error": "Hedef siteden yanıt okunurken hata oluştu."})
			return
		}

		// Hedef siteden gelen JSON yanıtını parse et
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Hedef siteden gelen yanıt parse edilemedi:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Hedef siteden gelen yanıt parse edilemedi."})
			return
		}

		// "word" alanını terminale yazdır
		if wordValue, ok := result["word"]; ok {
			fmt.Println("Bulunan word değeri:", wordValue)
		} else {
			fmt.Println("Yanıtta 'word' alanı bulunamadı.")
		}

		// Kullanıcıya basit bir yanıt döndür
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "İşlem başarılı, terminale yazdırıldı."})
	})

	// Sunucuyu başlat
	fmt.Println("Sunucu http://localhost:8080 adresinde çalışıyor...")
	http.ListenAndServe(":8080", nil)
}
