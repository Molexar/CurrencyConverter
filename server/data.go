package server

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	errorSleepTime   = time.Minute * 1
	successSleepTime = time.Hour * 1

	EcbCurrURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml" // Current url of xml data currencies
	currencyDateFormat = "2006-01-02" // The time format in ECB XML

	eur = "EUR"
)

//Starting gorountine what fetches currencies
func (s *Server) startCurrencyUpdating() {
	log.Println("Starting currency fetching")
	go func() {
		for {
			log.Println("Starting new currency fetching")

			// Initialize default naptime
			napTime := successSleepTime

			if data, errFetch := fetchCurrencyData(); errFetch == nil {
				if ts, curr, errParse := parseCurrencyData(data); errParse == nil {
					//updating currency data
					s.mutex.Lock()
					s.hasCurrencies, s.lastUpdateTime, s.currencies = true, ts, curr
					s.mutex.Unlock()

					log.Println("Currencies updated!")

					go s.callWebhooks()
				} else {
					log.Println("Error parsing data")
					napTime = errorSleepTime
				}
			} else {
				log.Println("Error fetching data")
				napTime = errorSleepTime
			}
			log.Println("Sleeping")
			time.Sleep(napTime)
		}
	}()
}

//Currency xml data
type Currency struct {
	Sender string `xml:"Sender>name"`
	Cube   []cube `xml:"Cube>Cube>Cube"`
}

//Parsed currency data
type cube struct {
	name string  `xml:"currency,attr"`
	rate float32 `xml:"rate, attr"`
}

// Time xml data
type timeCurrency struct{
	Time timeCube `xml:"Cube>Cube"`
}

type timeCube struct{
	time string `xml:"time,attr"`
}

func fetchCurrencyData() (data []byte, err error) {
	res, err := http.Get(EcbCurrURL)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return data, err
}

func parseCurrencyData(data []byte) (ts time.Time, currencies map[string]float32, err error) {
	//Parse to get Currency structured data
	var c Currency
	err = xml.Unmarshal(data, &c)
	if err != nil{
		return time.Now(), nil, err
	}

	//Parse again to get timeCurrency data
	var t timeCurrency
	err = xml.Unmarshal(data, &t)
	if (err != nil){
		return time.Now(), nil, err
	}
	
	ts, err = time.Parse(currencyDateFormat, t.Time.time)
	if err != nil{
		return time.Now(), nil, err
	}

	currencies = make(map[string]float32)

	currencies[eur] = 1

	//Insert all rates
	for _, cucurrency := range c.Cube{
		currencies[cucurrency.name] = cucurrency.rate
	}

	return ts, currencies, nil
}
