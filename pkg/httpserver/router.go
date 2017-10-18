package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kminehart/ambassador-auth-jwt/pkg/token"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/auth", token.DecodeHttpHandler).
		Methods(http.MethodPost)

	return r
}
