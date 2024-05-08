// server/http.go
package server

import (
	"PortalCRG/internal"
	"fmt"
	"net/http"
)

// HTTPServer representa el servidor HTTP.
type HTTPServer struct {
	UserService internal.UserService
}

// NewHTTPServer crea una nueva instancia de HTTPServer.
func NewHTTPServer(userService internal.UserService) *HTTPServer {
	return &HTTPServer{
		UserService: userService,
	}
}

// Start inicia el servidor HTTP.
func (s *HTTPServer) Start(port string) error {
	http.HandleFunc("/saludo", s.handleGreet)
	http.Handle("/", http.FileServer(http.Dir("./static/browser/")))
	http.HandleFunc("/portal/login", s.handleLogin)

	fmt.Printf("Servidor escuchando en el puerto %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func (s *HTTPServer) handleGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buscando index!!")
	fmt.Fprintf(w, s.UserService.Greet())
}

func (s *HTTPServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	alias := r.FormValue("alias")
	password := r.FormValue("password")

	user, err := s.UserService.AuthenticateUser(alias, password)
	if err != nil {
		http.Error(w, "Error de autenticación", http.StatusUnauthorized)
		return
	}

	// Autenticación exitosa
	// Aquí podrías manejar la sesión del usuario, por ejemplo, estableciendo una cookie de sesión.

	fmt.Fprintf(w, "¡Bienvenido, %s!", user.Name)
}
