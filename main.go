package main

import (
	"cmp"
	"fmt"
	"net/http"
	"os"
)

func main() {

	fs := http.FileServer(http.Dir("./static/browser/"))
	http.Handle("/", fs)

	http.HandleFunc("/saludo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Buscando index!!")
		fmt.Fprintf(w, "TODO OK!!")

	})
	fmt.Println("Servidor escuchando en el puerto :80, sin novedades")
	port := cmp.Or(os.Getenv("PORT"), "80")

	http.ListenAndServe(":"+port, nil)
}
