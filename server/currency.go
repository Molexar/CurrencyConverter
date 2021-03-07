package server

import "fmt"

func (s *Server) createResponse(base string) (r *currencyResponse, err error) {
	//Check if we has currencies
	if !s.hasCurrencies {
		return nil, fmt.Errorf("No currencies")
	}

	baserate, found := s.currencies[base]
	if !found {
		return nil, fmt.Errorf("Base string has no found")
	}

	//Create the response struct
	response := currencyResponse{}
	response.BaseCurrency = base
	response.CurrencyDate = s.lastUpdateTime.format("2006-01-02")

	// fill rates
	for name, rate := range s.currencies {
		relativeRate := rate / baserate
		r := rateResponse{
			Name: name,
			Rate: relativeRate,
		}
		response.Rates = append(response.Rates, r)
	}

	return &response, nil
}

func (s *Server) createConvertResponse()
