package httpapi

import "net/http"

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	s.sendData(w, r, http.StatusOK, map[string]any{"status": "ok"})
}
