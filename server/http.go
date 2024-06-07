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
	PortalService    internal.PortalRetroGamerService
	sessionsInServer map[string]string
}

func (s *HTTPServer) hashAlias(alias string) string {
	hasher := md5.New()
	hasher.Write([]byte(alias))
	return hex.EncodeToString(hasher.Sum(nil))
}

// NewHTTPServer crea una nueva instancia de HTTPServer.
func NewHTTPServer(portalRetroGamerService internal.PortalRetroGamerService) *HTTPServer {
	return &HTTPServer{
		PortalService: portalRetroGamerService,
	}
}

// Start inicia el servidor HTTP.
func (s *HTTPServer) Start(port string) error {

	s.sessionsInServer = make(map[string]string)

	http.HandleFunc("/saludo", s.handleGreet)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static/browser/"))))

	http.HandleFunc("/portal/login", s.handleLogin)

	http.HandleFunc("/portal/logout", s.handleLogout)

	http.HandleFunc("/portal/savePassword", s.savePassword)

	http.HandleFunc("/portal/saveProfile", s.savePerfil)

	http.HandleFunc("/portal/isOnline", s.isLogin)

	http.HandleFunc("/portal/userData", s.userData)

	http.HandleFunc("/portal/saveTips", s.saveTips)

	http.HandleFunc("/portal/deleteTips", s.deleteTip)

	http.HandleFunc("/portal/comment", s.comment)

	http.HandleFunc("/public/team", s.teams)

	http.HandleFunc("/public/checkCode", s.checkCode)

	http.HandleFunc("/public/checkAlias", s.checkAlias)

	http.HandleFunc("/public/userInfo", s.userInfo)

	http.HandleFunc("/public/createUser", s.createUser)

	http.HandleFunc("/public/tips", s.tips)

	http.HandleFunc("/public/new", s.getTips)

	http.HandleFunc("/s", s.sharedTips)

	http.HandleFunc("/public/loadTips", s.loadTips)

	http.HandleFunc("/public/loadTipsByPerfil", s.loadTipsPerfil)

	http.HandleFunc("/public/search", s.loadTipsSearch)

	fmt.Printf("Servidor escuchando en el puerto %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func (s *HTTPServer) handleGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buscando index!!")
	fmt.Fprintf(w, s.PortalService.Greet())
}

func (s *HTTPServer) MakeErrorMessage(w http.ResponseWriter, message string, code int) {
	error := ResponseMessage{}
	error.Code = code
	error.Message = message

	http.Error(w, message, http.StatusInternalServerError)
}
