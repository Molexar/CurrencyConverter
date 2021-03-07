package server

import (
	"encoding/json"
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

type currencyRequest struct {
	BaseCurrency string `json:"base_currency"`
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
		s.currencyHits.Add(1)
	} else if r.Request == http.MethodPost{
		var req currencyRequest
		err := s.getJsonRequest(r, &req)
		if err != nil{
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		res, err := s.createResponse(req.BaseCurrency)
		s.respondJson(w, res, err)
		s.currencyHits.Add(1)
	} else{
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (s *Server) respondJson(w http.ResponseWriter, v interface{}, err error){
	w.Header().Add("Content-Type", "application-json")
	if err != nil{
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	
	err = json.NewEncoder(w).Encode(v)

	if err != nil{
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (s *Server) getJsonRequest(r *http.Request, v interface{}) (err error){
	err = json.NewDecoder(r.Body).Decode(&v)
	if err != nil{
		return err
	}

	return nil
}
