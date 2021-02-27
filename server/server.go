package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	HostEnvironment = "GFS_CURRENCY_HOST" // hostname environment variable
	PortEnvironment = "GFS_CURRENCY_PORT" // environment portname

	defaultHost = "127.0.0.1" // default hostname
	defaultPort = "4000"
)

type webhook struct {
	BaseCurrency string `json:"base_currency"`
	URL          string `json:"url"`
	Secret       string `json:"secret"`
}

type Server struct {
	host string
	port int

	hasCurrencies  bool               // Информация, были ли получены данные валют
	lastUpdateTime time.Time          // Время последнего обновления курса валют
	currencies     map[string]float32 // курсы валют

	mutex    *sync.Mutex        // Мьютекс для синхронизации
	webhooks map[string]webhook // вебхуки

}

func New() (s *Server, err error) {
	// Initializing new Server object
	// Getting variables from environment
	host := getEnv(HostEnvironment, defaultHost)
	portStr := getEnv(PortEnvironment, defaultPort)
	port, err := strconv.Atoi(portStr)

	//Returning new Server object and error
	if err != nil {
		return nil, fmt.Errorf("Error parcing port number: %s", portStr)
	}
	return &Server{
		host: host,
		port: port,

		hasCurrencies: false,

		mutex:    &sync.Mutex{},
		webhooks: make(map[string]webhook),
	}, nil
}

//Runs server and returning errors
func (s *Server) Run() (err error) {
	log.Printf("Starting server on: Host=%s Port=%d\n", s.host, s.port)

	s.startCurrencyUpdating()
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), s)
}

func getEnv(env string, def string) (value string) {
	value = os.Getenv(env)
	if value == "" {
		value = def
	}
	return value
}
