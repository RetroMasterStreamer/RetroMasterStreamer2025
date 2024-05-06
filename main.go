package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	fmt.Println("Servidor escuchando en el puerto :80, no veo niun problema")
	http.ListenAndServe(":80", nil)
}
