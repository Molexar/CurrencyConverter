package server

import (
	"log"
	"net/http"
	"strings"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s]: %s\n", r.Method, r.RequestURI)

	//Check if has currencies
	if !s.hasCurrencies {
		log.Println("No currencies, error")
		http.Error(w, "No currencies", http.StatusServiceUnavailable)
		return
	}

	switch r.RequestURI {
	case "/currencies":
		s.currenciesHandler(w, r)
	case "/convert":
		s.convertHandler(w, r)
	case "/webhook":
		s.webhookHandler(w, r)
	default:
		if strings.HasPrefix(r.RequestURI, "/script?base=") {
			s.scriptHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}
