package main

import (
	"./broadcastd"
	"fmt"
	"log"
)

func main() {
	p := broadcastd.NewPubsub()
	s := broadcastd.NewServer(p)

	go p.Run()

	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("main: %v", err))
	}

	log.Printf("Listening on http://localhost%s\n", s.Addr)

	<-make(chan bool)
}
