// server/http.go
package server

import (
	"PortalCRG/internal/repository/entity"
	"PortalCRG/internal/util"
	"log"
	"strings"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid" // Para generar UUIDs únicos
)

func (s *HTTPServer) isLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("portal_ident")
	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessionsInServer[sessionToken]

	response := ResponseOnline{}

	userOnline, error := s.PortalService.GetStatusLogin(sessionToken, userID)

	if !ok || error != nil {
		// El token de sesión no coincide con ninguna sesión activa, redirigir al inicio de sesión
		response.Code = 401
		response.Status = "GAME OVER"
	}

	if userOnline != nil && userID == userOnline.Hash {
		response.Code = 200
		response.Status = "PLAYER"
		response.User = *userOnline
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (s *HTTPServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Configurar encabezados CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	if r.Method != http.MethodPost {
		s.MakeErrorMessage(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}

	// Definir estructura para datos de credenciales

	// Decodificar el cuerpo JSON en la estructura de credenciales
	var creds Credentials
	if err := json.Unmarshal(body, &creds); err != nil {
		s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
		return
	}

	if creds.Alias == "" || creds.Password == "" {
		s.MakeErrorMessage(w, "Como que te falto ingresar algo", http.StatusUnauthorized)
		return
	}

	// Autenticar usuario con las credenciales proporcionadas
	user, err := s.PortalService.AuthenticateUser(creds.Alias, creds.Password)
	if err != nil {
		s.MakeErrorMessage(w, "Error de autenticación", http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.New().String()

	response := ResponseLogin{}
	response.User = *user

	response.Hash = sessionToken

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	// Almacenar el token de sesión en el mapa de sesiones

	hash := s.hashAlias(response.User.Alias + r.RemoteAddr)
	s.sessionsInServer[sessionToken] = hash // Aquí puedes almacenar el ID de usuario u otra información relacionada con la sesión

	log.Println("Sesion TOKEN :" + sessionToken + "| User :" + creds.Alias + "| Hash :" + hash)

	// Establecer una cookie con el token de sesión
	http.SetCookie(w, &http.Cookie{
		Name:     "portal_ident",
		Value:    sessionToken,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
		// Otras configuraciones de cookie, como Path, MaxAge, etc.
	})

	s.PortalService.SetStatusLogin(creds.Alias, sessionToken, hash, true)

	w.Header().Set("Content-Type", "application/json")

	// Escribir la respuesta JSON en el cuerpo de la respuesta HTTP
	w.Write(jsonResponse)
}

func (s *HTTPServer) handleLogout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("portal_ident")
	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessionsInServer[sessionToken]
	if !ok {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	userOnline, err := s.PortalService.GetStatusLogin(sessionToken, userID)
	w.Header().Set("Content-Type", "application/json")

	if userOnline != nil {

		jsonResponse, err := s.logout(userOnline, sessionToken, userID)
		if err != nil {
			s.MakeErrorMessage(w, "Error al generar respuesta JSON"+err.Error(), http.StatusInternalServerError)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "portal_ident",
				Value:    "",
				SameSite: http.SameSiteNoneMode,
				// Otras configuraciones de cookie, como Path, MaxAge, etc.
			})
			// Escribir la respuesta JSON en el cuerpo de la respuesta HTTP
			w.Write(jsonResponse)
		}
	}
}

func (s *HTTPServer) logout(userOnline *entity.UserOnline, sessionToken string, userID string) ([]byte, error) {
	userData, err := s.PortalService.GetUserByAlias(userOnline.Alias)
	if err == nil {
		s.PortalService.SetStatusLogin(userData.Alias, sessionToken, userID, false)

		response := entity.UserOnline{}
		response.Alias = userData.Alias
		response.Online = false

		jsonResponse, _ := json.Marshal(response)

		return jsonResponse, nil
	} else {

		return nil, err
	}
}

func (s *HTTPServer) userData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("portal_ident")
	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessionsInServer[sessionToken]

	userOnline, error := s.PortalService.GetStatusLogin(sessionToken, userID)

	userData := entity.User{}
	jsonResponse, _ := json.Marshal(userData)

	if !ok || error != nil {

		if err != nil {
			s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
			return
		}

	} else {

		if userID == userOnline.Hash {
			/*
				http.SetCookie(w, &http.Cookie{
					Name:     "portal_ident",
					Value:    sessionToken,
					SameSite: http.SameSiteNoneMode,
					// Otras configuraciones de cookie, como Path, MaxAge, etc.
				})
			*/
			userData, err := s.PortalService.GetUserByAlias(userOnline.Alias)
			if err == nil {
				jsonResponse, err = json.Marshal(userData)
				if err != nil {
					s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
					return
				}

			} else {
				s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
				return
			}

		}
	}
	w.Write(jsonResponse)
}

func (s *HTTPServer) savePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("portal_ident")

	if r.Method != http.MethodPost {
		s.MakeErrorMessage(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessionsInServer[sessionToken]

	response := ResponseOnline{}
	response.Code = 500

	userOnline, error := s.PortalService.GetStatusLogin(sessionToken, userID)

	if !ok || error != nil {
		// El token de sesión no coincide con ninguna sesión activa, redirigir al inicio de sesión
		response.Code = 401
		response.Status = "GAME OVER"
	}

	if userID == userOnline.Hash {
		/*
			http.SetCookie(w, &http.Cookie{
				Name:     "portal_ident",
				Value:    sessionToken,
				SameSite: http.SameSiteNoneMode,
				// Otras configuraciones de cookie, como Path, MaxAge, etc.
			})
		*/
		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
			return
		}

		// Definir estructura para datos de credenciales

		// Decodificar el cuerpo JSON en la estructura de credenciales
		var passwordChange ChangePassword
		if err := json.Unmarshal(body, &passwordChange); err != nil {
			s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
			return
		}

		usuario, err := s.PortalService.GetUserByAlias(userOnline.Alias)
		if err != nil {
			s.MakeErrorMessage(w, "Error no existe el usuario", http.StatusInternalServerError)
			return
		}

		if passwordChange.NewPassword == passwordChange.ConfirmNewPassword && passwordChange.Password == usuario.Password {
			s.PortalService.ChangePassword(userOnline.Alias, passwordChange.NewPassword)
			response.Code = 200
		}

	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (s *HTTPServer) savePerfil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("portal_ident")

	if r.Method != http.MethodPost {
		s.MakeErrorMessage(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessionsInServer[sessionToken]

	response := ResponseOnline{}
	response.Code = 500

	userOnline, error := s.PortalService.GetStatusLogin(sessionToken, userID)

	if !ok || error != nil {
		// El token de sesión no coincide con ninguna sesión activa, redirigir al inicio de sesión
		response.Code = 401
		response.Status = "GAME OVER"
	}

	if userID == userOnline.Hash {
		/*
			http.SetCookie(w, &http.Cookie{
				Name:     "portal_ident",
				Value:    sessionToken,
				SameSite: http.SameSiteNoneMode,
				// Otras configuraciones de cookie, como Path, MaxAge, etc.
			})
		*/
		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
			return
		}

		// Definir estructura para datos de credenciales

		// Decodificar el cuerpo JSON en la estructura de credenciales
		var newUser entity.User
		if err := json.Unmarshal(body, &newUser); err != nil {
			s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
			return
		}

		usuario, err := s.PortalService.GetUserByAlias(userOnline.Alias)
		if err != nil {
			s.MakeErrorMessage(w, "Error no existe el usuario", http.StatusInternalServerError)
			return
		}

		usuarioRef, _ := s.PortalService.GetUserByTextRefer(newUser.ReferenceText)
		if usuarioRef != nil && usuarioRef.Alias != newUser.Alias {
			s.MakeErrorMessage(w, "Frase ya esta siendo ocupada", http.StatusInternalServerError)
			return
		}

		usuario.AboutMe = newUser.AboutMe
		usuario.Name = newUser.Name
		usuario.RRSS = newUser.RRSS

		for _, rrss := range newUser.RRSS {
			if rrss.Type == "youtube" && strings.Contains(rrss.URL, "youtube.com/") {
				usuario.AvatarYT = util.GetAvatarByURL(rrss.URL)
				break
			}
		}

		usuario.ReferenceText = newUser.ReferenceText

		s.PortalService.SaveUser(*usuario)
		response.Code = 200

	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
