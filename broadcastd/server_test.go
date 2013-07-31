package broadcastd

import (
	"github.com/bmizerany/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type ClosableRecorder struct {
	*httptest.ResponseRecorder
	closer chan bool
}

func NewClosableRecorder() *ClosableRecorder {
	r := httptest.NewRecorder()
	closer := make(chan bool)
	return &ClosableRecorder{r, closer}
}

func (r *ClosableRecorder) CloseNotify() <-chan bool {
	return r.closer
}

func TestServer(t *testing.T) {
	p := NewPubsub()
	s := NewServer(p)

	go p.Run()

	// Subscribe
	get, _ := http.NewRequest("GET", "/", nil)
	sub := NewClosableRecorder()
	go s.ServeHTTP(sub, get)

	time.Sleep(time.Second)

	// Publish
	body := strings.NewReader("foo")
	post, _ := http.NewRequest("POST", "/", body)
	s.ServeHTTP(httptest.NewRecorder(), post)
	body = strings.NewReader("bar")
	post, _ = http.NewRequest("POST", "/", body)
	s.ServeHTTP(httptest.NewRecorder(), post)

	time.Sleep(time.Second)
	sub.closer <- true

	assert.Equal(t, sub.Body.String(), "foo\nbar\n")
}
