// server/http.go
package server

import (
	"PortalCRG/internal/repository/entity"
	"PortalCRG/internal/util"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid" // Para generar UUIDs únicos
)

func (s *HTTPServer) isLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessions[sessionToken]

	response := ResponseOnline{}

	userOnline, error := s.UserService.GetStatusLogin(sessionToken, userID)

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
		http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
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
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}

	// Definir estructura para datos de credenciales

	// Decodificar el cuerpo JSON en la estructura de credenciales
	var creds Credentials
	if err := json.Unmarshal(body, &creds); err != nil {
		http.Error(w, "Formato de datos incorrecto", http.StatusBadRequest)
		return
	}

	// Autenticar usuario con las credenciales proporcionadas
	user, err := s.UserService.AuthenticateUser(creds.Alias, creds.Password)
	if err != nil {
		http.Error(w, "Error de autenticación", http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.New().String()

	response := ResponseLogin{}
	response.User = *user

	response.Hash = sessionToken

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	s.sessions = make(map[string]string)
	// Almacenar el token de sesión en el mapa de sesiones
	hash := s.hashAlias(response.User.Alias)
	s.sessions[sessionToken] = hash // Aquí puedes almacenar el ID de usuario u otra información relacionada con la sesión

	// Establecer una cookie con el token de sesión
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		// Otras configuraciones de cookie, como Path, MaxAge, etc.
	})

	s.UserService.SetStatusLogin(creds.Alias, sessionToken, hash, true)

	w.Header().Set("Content-Type", "application/json")

	// Escribir la respuesta JSON en el cuerpo de la respuesta HTTP
	w.Write(jsonResponse)
}

func (s *HTTPServer) handleLogout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessions[sessionToken]
	if !ok {
		http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	userOnline, err := s.UserService.GetStatusLogin(sessionToken, userID)
	w.Header().Set("Content-Type", "application/json")

	if userOnline != nil {

		jsonResponse, err := s.logout(userOnline, sessionToken, userID)
		if err != nil {
			http.Error(w, "Error al generar respuesta JSON"+err.Error(), http.StatusInternalServerError)
		} else {

			// Escribir la respuesta JSON en el cuerpo de la respuesta HTTP
			w.Write(jsonResponse)
		}
	}
}

func (s *HTTPServer) logout(userOnline *entity.UserOnline, sessionToken string, userID string) ([]byte, error) {
	userData, err := s.UserService.GetUserByAlias(userOnline.Alias)
	if err == nil {
		s.UserService.SetStatusLogin(userData.Alias, sessionToken, userID, false)

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
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessions[sessionToken]

	userOnline, error := s.UserService.GetStatusLogin(sessionToken, userID)

	userData := entity.User{}
	jsonResponse, _ := json.Marshal(userData)

	if !ok || error != nil {

		if err != nil {
			http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
			return
		}

	} else {

		if userID == userOnline.Hash {
			userData, err := s.UserService.GetUserByAlias(userOnline.Alias)
			if err == nil {
				jsonResponse, err = json.Marshal(userData)
				if err != nil {
					http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
					return
				}

			} else {
				http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
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
	cookie, err := r.Cookie("session_token")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessions[sessionToken]

	response := ResponseOnline{}
	response.Code = 500

	userOnline, error := s.UserService.GetStatusLogin(sessionToken, userID)

	if !ok || error != nil {
		// El token de sesión no coincide con ninguna sesión activa, redirigir al inicio de sesión
		response.Code = 401
		response.Status = "GAME OVER"
	}

	if userID == userOnline.Hash {

		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
			return
		}

		// Definir estructura para datos de credenciales

		// Decodificar el cuerpo JSON en la estructura de credenciales
		var passwordChange ChangePassword
		if err := json.Unmarshal(body, &passwordChange); err != nil {
			http.Error(w, "Formato de datos incorrecto", http.StatusBadRequest)
			return
		}

		usuario, err := s.UserService.GetUserByAlias(userOnline.Alias)
		if err != nil {
			http.Error(w, "Error no existe el usuario", http.StatusInternalServerError)
			return
		}

		if passwordChange.NewPassword == passwordChange.ConfirmNewPassword && passwordChange.Password == usuario.Password {
			s.UserService.ChangePassword(userOnline.Alias, passwordChange.NewPassword)
			response.Code = 200
		}

	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (s *HTTPServer) savePerfil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	cookie, err := r.Cookie("session_token")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		// El token de sesión no está presente o es inválido, redirigir al inicio de sesión
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verificar el token de sesión en el mapa de sesiones
	sessionToken := cookie.Value
	userID, ok := s.sessions[sessionToken]

	response := ResponseOnline{}
	response.Code = 500

	userOnline, error := s.UserService.GetStatusLogin(sessionToken, userID)

	if !ok || error != nil {
		// El token de sesión no coincide con ninguna sesión activa, redirigir al inicio de sesión
		response.Code = 401
		response.Status = "GAME OVER"
	}

	if userID == userOnline.Hash {

		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
			return
		}

		// Definir estructura para datos de credenciales

		// Decodificar el cuerpo JSON en la estructura de credenciales
		var newUser entity.User
		if err := json.Unmarshal(body, &newUser); err != nil {
			http.Error(w, "Formato de datos incorrecto", http.StatusBadRequest)
			return
		}

		usuario, err := s.UserService.GetUserByAlias(userOnline.Alias)
		if err != nil {
			http.Error(w, "Error no existe el usuario", http.StatusInternalServerError)
			return
		}

		usuario.AboutMe = newUser.AboutMe
		usuario.Name = newUser.Name
		usuario.RRSS = newUser.RRSS

		youtubeURL := ""
		for _, rrss := range newUser.RRSS {
			if rrss.Type == "youtube" {
				youtubeURL = rrss.URL
			}
		}

		usuario.AvatarYT = util.GetAvatarByURL(youtubeURL)

		usuario.ReferenceText = newUser.ReferenceText

		s.UserService.SaveUser(*usuario)
		response.Code = 200

	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
