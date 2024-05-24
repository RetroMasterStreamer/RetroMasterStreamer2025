package server

import (
	"PortalCRG/internal/repository/entity"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	// Para generar UUIDs únicos
)

func (s *HTTPServer) checkAlias(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}

	var user entity.User
	if err := json.Unmarshal(body, &user); err != nil {
		s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
		return
	}

	userRef, err := s.PortalService.GetUserByAlias(user.Alias)
	if err != nil {
		userRef = &user
	}

	if userRef == nil {
		userRef = &user
	}

	jsonResponse, err := json.Marshal(userRef)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func (s *HTTPServer) checkCode(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}

	var code Credentials
	if err := json.Unmarshal(body, &code); err != nil {
		s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
		return
	}

	iddqd := code.Password

	if iddqd == "" {
		s.MakeErrorMessage(w, "Que intentas hacer?", http.StatusInternalServerError)
		return
	}

	userRef, err := s.PortalService.GetUserByRefer(iddqd)
	if err != nil {
		s.MakeErrorMessage(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(userRef)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func (s *HTTPServer) teams(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	teams, err := s.PortalService.GetAllUsers()
	if err != nil {
		s.MakeErrorMessage(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(teams)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func (s *HTTPServer) createUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}

	var createUserRequest NewUserRequest
	if err := json.Unmarshal(body, &createUserRequest); err != nil {
		s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
		return
	}

	userByCode, err := s.PortalService.GetUserByRefer(createUserRequest.Code)
	if err != nil {
		s.MakeErrorMessage(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}

	if userByCode.Alias == createUserRequest.RefUser.Alias {
		createUserRequest.NewUser.UserRef = userByCode.Alias
		newUser := entity.User{}
		newUser.Alias = createUserRequest.NewUser.Alias
		newUser.Password = createUserRequest.NewUser.Password
		newUser.UserRef = createUserRequest.RefUser.Alias
		errCreate := s.PortalService.CreateUser(&newUser)
		if errCreate != nil {
			s.MakeErrorMessage(w, "Error al crear "+errCreate.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse, err := json.Marshal(createUserRequest.NewUser)
		if err != nil {
			s.MakeErrorMessage(w, "Error no puedo crear el dato", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	} else {
		s.MakeErrorMessage(w, "Error no puedo crear el dato", http.StatusInternalServerError)
		return
	}

}

func (s *HTTPServer) saveTips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}

	// Definir estructura para datos de credenciales

	// Decodificar el cuerpo JSON en la estructura de credenciales
	var tips entity.PostNew
	if err := json.Unmarshal(body, &tips); err != nil {
		s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
		return
	}

	s.PortalService.CreateTips(&tips)

	return

}

func (s *HTTPServer) loadTipsPerfil(w http.ResponseWriter, r *http.Request) {
	alias := r.URL.Query().Get("alias")
	if alias == "" {
		http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	skip := int64(page * limit)

	tips, err := s.PortalService.GetTipsByAliasWithPagination(alias, skip, int64(limit))
	if err != nil {
		http.Error(w, "Error fetching tips", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tips); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) loadTips(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	skip := int64(page * limit)

	tips, err := s.PortalService.GetTipsWithPagination(skip, int64(limit))
	if err != nil {
		http.Error(w, "Error fetching tips", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tips); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
func (s *HTTPServer) tips(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	teams, err := s.PortalService.GetAllTips()
	if err != nil {
		s.MakeErrorMessage(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(teams)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func (s *HTTPServer) getTips(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}

	tip := s.PortalService.GetTipByID(id)

	jsonResponse, err := json.Marshal(tip)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func (s *HTTPServer) loadTipsSearch(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	skip := int64(page * limit)

	tips, err := s.PortalService.GetTipsWithSearch(search, skip, int64(limit))
	if err != nil {
		s.MakeErrorMessage(w, "No existen resultados :(", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tips); err != nil {
		s.MakeErrorMessage(w, "No pude crear una respuesta", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) deleteTip(w http.ResponseWriter, r *http.Request) {

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
	userID, ok := s.sessions[sessionToken]

	response := ResponseOnline{}

	userOnline, error := s.PortalService.GetStatusLogin(sessionToken, userID)

	if !ok || error != nil {
		// El token de sesión no coincide con ninguna sesión activa, redirigir al inicio de sesión
		response.Code = 401
		response.Status = "GAME OVER"
		s.MakeErrorMessage(w, "Intente logearse otra vez", http.StatusMethodNotAllowed)
	}

	if userOnline != nil && userID == userOnline.Hash {

		if r.Method != http.MethodDelete {
			s.MakeErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			ID string `json:"id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.MakeErrorMessage(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := s.PortalService.DeleteTip(req.ID, userOnline.Alias); err != nil {
			s.MakeErrorMessage(w, "Failed to delete tip", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *HTTPServer) userInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.MakeErrorMessage(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}

	var user entity.User
	if err := json.Unmarshal(body, &user); err != nil {
		s.MakeErrorMessage(w, "Formato de datos incorrecto", http.StatusBadRequest)
		return
	}

	userRef, err := s.PortalService.GetUserByAlias(user.Alias)
	if err != nil {
		userRef = &user
	}

	if userRef == nil {
		userRef = &user
	}

	jsonResponse, err := json.Marshal(userRef)
	if err != nil {
		s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func (s *HTTPServer) sharedTips(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	tip := s.PortalService.GetTipByID(id)
	tips := TipsShared{}
	if tip == nil {
		tips = TipsShared{
			Title:    "",
			ID:       "",
			URL:      "/",
			AvatarYT: "",
		}
	} else {

		author, _ := s.PortalService.GetUserByAlias(tip.Author)

		// Define your tips data
		tips = TipsShared{
			Title:    tip.Title,
			ID:       tip.ID,
			URL:      tip.URL,
			AvatarYT: author.AvatarYT,
		}
	}
	// Define the path to the HTML file
	htmlFilePath := filepath.Join("static", "browser", "shared.html")

	// Parse the HTML file as a template
	tmpl, err := template.ParseFiles(htmlFilePath)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Execute the template with the tips data and write to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, tips)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}
