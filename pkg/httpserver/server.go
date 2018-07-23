package httpserver

import (
	"fmt"
	"net/http"

	"github.com/tomwganem/ambassador-auth-jwt/pkg/token"
)

// Server needs to know about the Issuer url to verify tokens against
type Server struct {
	Issuer string
}

// Start accepting requests and decoding Authorization headers
func (Start *Server) Start(port int) error {
	http.HandleFunc("/", token.DecodeHTTPHandler)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}

// NewServer creates a new "Server" object
func NewServer(issuer string) *Server {
	return &Server{
		Issuer: issuer,
	}
}
