package broadcastd

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

type Server struct {
	Addr   string
	Config *Config
	Router *mux.Router
	pubsub *Pubsub
}

func NewServer(config *Config, pubsub *Pubsub) *Server {
	router := mux.NewRouter()
	addr := fmt.Sprintf(":%d", config.Port)

	s := &Server{addr, config, router, pubsub}

	s.HandleFunc("/", s.PublishHandler).Methods("POST")
	s.HandleFunc("/", s.SubscribeHandler).Methods("GET")

	return s
}

func (s *Server) HandleFunc(route string, fn func(http.ResponseWriter, *http.Request)) *mux.Route {
	return s.Router.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		if ok := s.BeforeHandler(w, req); ok {
			fn(w, req)
		}
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(w, req)
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)

	if err != nil {
		return err
	}

	srv := http.Server{Handler: s}
	go srv.Serve(listener)

	return nil
}
