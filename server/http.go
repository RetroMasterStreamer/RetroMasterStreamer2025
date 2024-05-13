// server/http.go
package server

import (
	"PortalCRG/internal"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	// Para generar UUIDs únicos
)

// HTTPServer representa el servidor HTTP.
type HTTPServer struct {
	UserService internal.UserService
	sessions    map[string]string
}

func (s *HTTPServer) hashAlias(alias string) string {
	hasher := md5.New()
	hasher.Write([]byte(alias))
	return hex.EncodeToString(hasher.Sum(nil))
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

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static/browser/"))))

	http.HandleFunc("/portal/login", s.handleLogin)

	http.HandleFunc("/portal/logout", s.handleLogout)

	http.HandleFunc("/portal/savePassword", s.savePassword)

	http.HandleFunc("/portal/saveProfile", s.savePerfil)

	http.HandleFunc("/portal/isOnline", s.isLogin)

	http.HandleFunc("/portal/userData", s.userData)

	fmt.Printf("Servidor escuchando en el puerto %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func (s *HTTPServer) handleGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buscando index!!")
	fmt.Fprintf(w, s.UserService.Greet())
}
