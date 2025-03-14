package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTPS over LAN!"))
	})

	log.Println("Server running at https://192.168.1.100:8443")

	err := http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}
