package main

import (
	"encoding/json"
	"net/http"
)

// query components are often used to carry identifying information in the form of "key=value" pairs
func anonimfunc(wrte http.ResponseWriter, req *http.Request) {
	value := req.URL.Query().Get("word")
	//Query parametre almıyor
	object := map[string]string{"word": value}
	wrte.Header()
	//header veya body
	json.NewEncoder(wrte).Encode(object) //bi üsttekini tutuyor
	// Encode(object) objecti jsona dönüştürüyor? zaten objectle aynı formatta??
}

func main() {
	http.HandleFunc("/", anonimfunc)
	//http.ResponseWriter, *http.Request kendisi alıyor, parametre girmedim
	//gelen tüm bilgieleri gönderiyouır
	//http.HandleFunc parametreleri otomatik olarak sağlıyormuş
	http.ListenAndServe(":8081", nil)
	//http://localhost:8081/?word=berfin
}
