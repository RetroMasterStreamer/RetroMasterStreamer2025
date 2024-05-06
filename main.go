package main

import (
	"cmp"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Contact struct {
	ID    int
	Name  string
	Phone string
}

func main() {
	// Abrir la base de datos SQLite
	db, err := sql.Open("sqlite3", "./agenda.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear la tabla si no existe
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			phone TEXT
		)`)
	if err != nil {
		log.Fatal(err)
	}

	// Rutas para manejar las solicitudes HTTP
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		// Manejar la creación de un nuevo contacto
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		stmt, err := db.Prepare("INSERT INTO contacts(name, phone) VALUES(?,?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(name, phone)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Contacto creado correctamente")
	})

	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		// Manejar la lectura de todos los contactos
		rows, err := db.Query("SELECT * FROM contacts")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var contacts []Contact
		for rows.Next() {
			var contact Contact
			err := rows.Scan(&contact.ID, &contact.Name, &contact.Phone)
			if err != nil {
				log.Fatal(err)
			}
			contacts = append(contacts, contact)
		}

		for _, contact := range contacts {
			fmt.Fprintf(w, "ID: %d, Name: %s, Phone: %s\n", contact.ID, contact.Name, contact.Phone)
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		// Manejar la actualización de un contacto existente
		id := r.FormValue("id")
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		stmt, err := db.Prepare("UPDATE contacts SET name=?, phone=? WHERE id=?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(name, phone, id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Contacto actualizado correctamente")
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// Manejar la eliminación de un contacto
		id := r.FormValue("id")
		stmt, err := db.Prepare("DELETE FROM contacts WHERE id=?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Contacto eliminado correctamente")
	})

	// Iniciar el servidor

	fmt.Println("Servidor escuchando en el puerto :80, sin novedades")
	port := cmp.Or(os.Getenv("PORT"), "80")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Buscando index!!")
		http.ServeFile(w, r, "static/index.html")
	})

	http.ListenAndServe(":"+port, nil)
}
