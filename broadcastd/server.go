package broadcastd

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type Server struct {
	Addr   string
	Config *Config
	pubsub *Pubsub
}

func NewServer(config *Config, pubsub *Pubsub) *Server {
	addr := fmt.Sprintf(":%d", config.Port)
	return &Server{addr, config, pubsub}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		c := s.pubsub.Subscribe()

		go c.Write(w)

		if close, ok := w.(http.CloseNotifier); ok {
			<-close.CloseNotify()
			s.pubsub.Unsubscribe(c)
		}
	} else if req.Method == "POST" {
		value, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}

		s.pubsub.Publish(string(value))
	}
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
