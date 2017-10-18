package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Secret string
}

func (this *Server) Start(port int) error {
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), this.Router)
}

func NewServer(secret string) *Server {
	router := GetRouter()
	return &Server{
		Secret: secret,
		Router: router,
	}
}
