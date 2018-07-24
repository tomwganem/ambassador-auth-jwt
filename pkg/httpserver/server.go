package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tomwganem/ambassador-auth-jwt/pkg/token"
	"gopkg.in/square/go-jose.v2"
)

var (
	// JwtCheckExp will determine if we need to verify if the token is expired or not
	JwtCheckExp = true
	// JwtIssuer is the url where we can retreive a set of public keys to verify rsa based tokens with
	JwtIssuer = "http://localhost/.well-known/jwks.json"
	// JwtOutboundHeader is the name of header the parsed token cliams will be inserted into
	JwtOutboundHeader = "X-JWT-PAYLOAD"
)

// Server needs to know about the Issuer url to verify tokens against
type Server struct {
	Issuer string
	JwkSet jose.JSONWebKeySet
}

// Start accepting requests and decoding Authorization headers
func (server *Server) Start(port int) error {
	// jwks, err := token.JwkSetGet(server.Issuer)
	http.HandleFunc("/", server.DecodeHTTPHandler)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}

// DecodeHTTPHandler will try to extract the bearer token found in the Authorization header of each request
func (server *Server) DecodeHTTPHandler(w http.ResponseWriter, r *http.Request) {
	errorLogger := log.WithFields(log.Fields{
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
		"method":      r.Method,
		"request_uri": r.RequestURI,
		"user_agent":  r.UserAgent(),
		"status":      "401",
	})
	// Get the Jwt
	jwt := r.Header.Get("Authorization")
	if jwt == "" {
		errorLogger.Error("Unable to retrieve Jwt from Authorization header")
		http.Error(w, fmt.Sprintf("Unable to retrieve Jwt from header '%s'", "Authorization"), 401)
		return
	}

	// Decode it
	mapClaims := make(map[string]interface{})
	jwt = strings.Replace(jwt, "Bearer ", "", 1)
	decoded, jwkset, err := token.Decode(jwt, server.JwkSet, JwtIssuer)
	server.JwkSet = jwkset
	if err != nil {
		errorLogger.Error(err.Error())
		http.Error(w, err.Error(), 401)
		return
	}
	claims, err := json.Marshal(decoded)
	if err := json.Unmarshal(claims, &mapClaims); err != nil {
		errorLogger.Error(err.Error())
		http.Error(w, err.Error(), 401)
		return
	}

	if JwtCheckExp {
		// Make sure the exp is before today...
		if _, ok := mapClaims["exp"]; ok != true {
			errorLogger.Error(err.Error())
			http.Error(w, err.Error(), 401)
			return
		}

		exp := time.Unix(mapClaims["exp"].(int64), 0)
		now := time.Now()

		if exp.Before(now) {
			errorLogger.Error("Token is expired")
			http.Error(w, "This token is expired", 401)
			return
		}
	}
	h := make(map[string]interface{})
	h[JwtOutboundHeader] = mapClaims
	log.WithFields(log.Fields{
		"remote_addr":     r.RemoteAddr,
		"host":            r.Host,
		"method":          r.Method,
		"request_uri":     r.RequestURI,
		"user_agent":      r.UserAgent(),
		"outbound_header": h,
		"status":          "200",
	}).Info("Authentication Success")
	w.Header().Set(JwtOutboundHeader, string(claims))
}

// NewServer creates a new "Server" object
func NewServer(issuer string) *Server {
	jwks, err := token.JwkSetGet(issuer)
	if err != nil {
		log.WithFields(log.Fields{
			"keyset": jwks,
			"issuer": issuer,
		}).Fatal("Unable to retreive keyset")
	}

	return &Server{
		Issuer: issuer,
		JwkSet: jwks,
	}
}
