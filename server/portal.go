package server

import (
	"PortalCRG/internal/repository/entity"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
