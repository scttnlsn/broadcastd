package broadcastd

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *Server) BeforeHandler(w http.ResponseWriter, req *http.Request) bool {
	if ok := auth(req, s.Config); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		send(w, Json{"error": "Unauthorized"})
		return false
	}

	return true
}

func (s *Server) PublishHandler(w http.ResponseWriter, req *http.Request) {
	value, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	s.pubsub.Publish(string(value))

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) SubscribeHandler(w http.ResponseWriter, req *http.Request) {
	c := s.pubsub.Subscribe()

	go c.Write(w)

	if close, ok := w.(http.CloseNotifier); ok {
		<-close.CloseNotify()
		s.pubsub.Unsubscribe(c)
	}
}

// Helpers

type Json map[string]interface{}

func send(w http.ResponseWriter, data Json) error {
	bytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)

	return nil
}

func auth(req *http.Request, config *Config) bool {
	if config.Auth == "" {
		return true
	}

	s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return false
	}

	base, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(base), ":", 2)
	if len(pair) != 2 {
		return false
	}

	password := pair[1]
	if config.Auth != password {
		return false
	}

	return true
}
