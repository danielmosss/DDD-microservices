package main

import (
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

	server := http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	return server.ListenAndServe()
}
