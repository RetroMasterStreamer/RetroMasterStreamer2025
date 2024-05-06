package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Buscando index!!")
		http.ServeFile(w, r, "static/index.html")
	})

	fmt.Println("Servidor escuchando en el puerto :80, sin novedades")
	http.ListenAndServe(":80", nil)
}
