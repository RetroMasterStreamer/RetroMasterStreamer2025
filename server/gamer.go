package server

import (
	"PortalCRG/internal/games"
	"net/http"
	"strings"
)

func (s *HTTPServer) PlayTheGame(w http.ResponseWriter, r *http.Request) {
	params := r.RequestURI[len("/public/games/"):]
	parts := strings.SplitN(params, "/", 2)
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	emulator, name := parts[0], parts[1]

	data, err := games.InsertCoin(emulator, name)
	if err != nil {
		s.MakeErrorMessage(w, "Error Catrige Malo", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
