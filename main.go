package main

import (
	"cmp"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Buscando index!!")
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/saludo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Buscando index!!")
		fmt.Fprintf(w, "TODO OK!!")

	})
	fmt.Println("Servidor escuchando en el puerto :80, sin novedades")
	port := cmp.Or(os.Getenv("PORT"), "80")

	http.ListenAndServe(":"+port, nil)
}
