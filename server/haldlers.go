package server

import (
	"log"
	"net/http"
	"strings"
)

type currencyResponse struct {
	CurrencyDate string         `json:"currency_date"`
	BaseCurrency string         `json:"base_currency"`
	Rates        []rateResponse `json:"rates"`
}

type rateResponse struct {
	Name string  `json:"name"`
	Rate float32 `json:"rate"`
}

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

//Handles request /currencies
func (s *Server) currenciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Create response with EUR base
		res, err := s.createResponse(eur)
		s.respondJson(w, res, err)
	}
}
