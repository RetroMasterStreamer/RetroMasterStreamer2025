package server

import (
	"PortalCRG/internal/repository/entity"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	// Para generar UUIDs únicos
)

func (s *HTTPServer) checkAlias(w http.ResponseWriter, r *http.Request) {
	log.Println("checkAlias  ")

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
	log.Println("checkCode  ")

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
	log.Println("teams  ")

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
	log.Println("createUser  ")

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
	log.Println("saveTips  ")
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

		tips.Author = userOnline.Alias

		errorNew, id := s.PortalService.CreateTips(&tips)

		if tips.Type == "download" {
			tips.File.IdTips = id
			erroInsertFile := s.DriveService.SaveFile(tips.File)

			if erroInsertFile != nil {
				log.Println("Eliminando tips dado que archivo esta malo")
				s.PortalService.DeleteTip(userOnline.Alias, tips.ID)

				s.MakeErrorMessage(w, "Error al subir archivos, estamos muy mal, avisale al administrador ayuda!!", http.StatusBadRequest)
				return
			}
		}

		if errorNew != nil {
			fmt.Println("Atencion : " + errorNew.Error())
		}

	}

}

func (s *HTTPServer) loadTipsPerfil(w http.ResponseWriter, r *http.Request) {
	log.Println("loadTipsPerfil  ")
	alias := r.URL.Query().Get("alias")
	if alias == "" {
		log.Println("Error ")
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
		log.Println("Error ")
		http.Error(w, "Error fetching tips", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tips); err != nil {
		log.Println("Error ")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) download(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")

	file, err := s.DriveService.GetFileByID(fileID)
	if err != nil {
		log.Println("Error ")
		http.Error(w, "Archivo no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
	w.Header().Set("Content-Type", file.Type)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))

	w.Write([]byte(file.Content))
}

func (s *HTTPServer) loadTips(w http.ResponseWriter, r *http.Request) {
	log.Println("loadTips  ")

	typeOfTips := []string{"tips", "youtube", "sitios"}

	typeOfTipsHeader := r.Header.Get("typeOfTips")
	if typeOfTipsHeader == "" {
		log.Println("No 'typeOfTips' header provided")
	} else {
		log.Printf("Received typeOfTips header: %s\n", typeOfTipsHeader)

		// Suponiendo que el array se envía como una cadena separada por comas: "videos,sitios,tips"
		typeOfTips = strings.Split(typeOfTipsHeader, ",")
		log.Printf("Parsed typeOfTips: %v\n", typeOfTips)

		// Aquí puedes usar 'typeOfTips' como quieras (por ejemplo, para filtrar los tips)
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

	tips, err := s.PortalService.GetTipsWithPagination(skip, int64(limit), typeOfTips)
	if err != nil {
		log.Println("Error ")
		http.Error(w, "Error fetching tips", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tips); err != nil {
		log.Println("Error ")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
func (s *HTTPServer) tips(w http.ResponseWriter, r *http.Request) {
	log.Println("tips  ")

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
	log.Println("getTips  ")
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("Error ")
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

	log.Println("loadTipsSearch  ")

	typeOfTips := []string{"tips", "youtube", "sitios"}

	typeOfTipsHeader := r.Header.Get("typeOfTips")
	if typeOfTipsHeader == "" {
		log.Println("No 'typeOfTips' header provided")
	} else {
		log.Printf("Received typeOfTips header: %s\n", typeOfTipsHeader)

		// Suponiendo que el array se envía como una cadena separada por comas: "videos,sitios,tips"
		typeOfTips = strings.Split(typeOfTipsHeader, ",")
		log.Printf("Parsed typeOfTips: %v\n", typeOfTips)

		// Aquí puedes usar 'typeOfTips' como quieras (por ejemplo, para filtrar los tips)
	}

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

	if contains(typeOfTips, "youtube") {
		s.PortalService.UpdateVideosTeams(search)
	}

	tips, err := s.PortalService.GetTipsWithSearch(search, skip, int64(limit), typeOfTips)
	if err != nil {
		s.MakeErrorMessage(w, "No existen resultados :(", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tips); err != nil {
		s.MakeErrorMessage(w, "No pude crear una respuesta", http.StatusInternalServerError)
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (s *HTTPServer) deleteTip(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteTip  ")

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
		} else {
			s.DriveService.DeleteFile(req.ID)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *HTTPServer) userInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("userInfo  ")

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
	log.Println("sharedTips !!! ")
	id := r.URL.Query().Get("id")
	tip := s.PortalService.GetTipByID(id)
	tips := TipsShared{}
	if tip == nil {
		tips = TipsShared{
			Title:    "",
			ID:       "",
			URL:      "/",
			Content:  "Portal Retro Gamer",
			AvatarYT: "",
		}
	} else {

		author, _ := s.PortalService.GetUserByAlias(tip.Author)

		// Define your tips data
		tips = TipsShared{
			Title:    tip.Title,
			ID:       tip.ID,
			URL:      tip.URL,
			Content:  tip.Content,
			AvatarYT: author.AvatarYT,
		}
	}
	// Define the path to the HTML file
	htmlFilePath := filepath.Join("static", "browser", "shared.html")

	// Parse the HTML file as a template
	tmpl, err := template.ParseFiles(htmlFilePath)
	if err != nil {
		log.Println("Error ")
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Execute the template with the tips data and write to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, tips)
	if err != nil {
		log.Println("Error ")
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) comment(w http.ResponseWriter, r *http.Request) {
	log.Println("comment  ")

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
		s.MakeErrorMessage(w, "Intente logearse otra vez", http.StatusMethodNotAllowed)
	}

	if userOnline != nil && userID == userOnline.Hash {

		comment := CommentPortal{}

		comment.Author = userOnline.Alias

		if r.Method != http.MethodPost {
			s.MakeErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			s.MakeErrorMessage(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if strings.Trim(comment.Comment, "") != "" {

			tipRetro := s.PortalService.GetTipByID(comment.ID)

			author, _ := s.PortalService.GetUserByAlias(comment.Author)

			commentRetro := entity.CommentRetro{}
			commentRetro.Author = comment.Author
			commentRetro.Comment = comment.Comment
			commentRetro.Avatar = author.AvatarYT
			commentRetro.Date = comment.Date

			tipRetro.Comments = append(tipRetro.Comments, commentRetro)

			tipRetro.Date = comment.Date

			err, _ := s.PortalService.CreateTips(tipRetro)
			if err != nil {
				s.MakeErrorMessage(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			jsonResponse, err := json.Marshal(tipRetro.Comments)
			if err != nil {
				s.MakeErrorMessage(w, "Error al generar respuesta JSON", http.StatusInternalServerError)
				return
			}
			w.Write(jsonResponse)
		} else {
			s.MakeErrorMessage(w, "Sin comentarios", http.StatusBadRequest)
		}

	}
}
