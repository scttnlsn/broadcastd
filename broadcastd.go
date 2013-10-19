package main

import (
	"flag"
	"fmt"
	"github.com/scttnlsn/broadcastd/broadcastd"
	"log"
	"runtime"
)

var config *broadcastd.Config

func init() {
	config = broadcastd.NewConfig()

	flag.UintVar(&config.Port, "port", 5454, "port on which to listen")
	flag.StringVar(&config.Auth, "auth", "", "HTTP basic auth password required for all requests")
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	p := broadcastd.NewPubsub()
	s := broadcastd.NewServer(config, p)

	go p.Run()

	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("main: %v", err))
	}

	log.Printf("Listening on http://localhost%s\n", s.Addr)

	<-make(chan bool)
}
