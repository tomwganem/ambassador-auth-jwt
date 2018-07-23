package httpserver

import (
	"fmt"
	"net/http"

	"../token"
)

type Server struct {
	Secret string
}

func (this *Server) Start(port int) error {
	http.HandleFunc("/", token.DecodeHTTPHandler)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}

func NewServer(secret string) *Server {
	return &Server{
		Secret: secret,
	}
}
