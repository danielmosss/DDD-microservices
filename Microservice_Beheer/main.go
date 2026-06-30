package main

import (
	//"beheer/database"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	mux := http.NewServeMux()

	//mux.HandleFunc{"GET/assets"}

	server := http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	log.Printf("Starting API server on %s", s.addr)

	return server.ListenAndServe()
}

func main() {
	//database.Init()

	server := NewAPIServer(":3000")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
