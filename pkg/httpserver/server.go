package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
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

// DecodeHTTPHandler will try to extract the bearer token found in the Authorization header of each request and verify it
func (server *Server) DecodeHTTPHandler(w http.ResponseWriter, r *http.Request) {
	errorLogger := log.WithFields(log.Fields{
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
		"url":         r.URL,
		"method":      r.Method,
		"request_uri": r.RequestURI,
		"user_agent":  r.UserAgent(),
		"status":      "401",
	})

	unauthorized := map[string]string{"code": "unauthorized", "message": "You are not authorized to perform the requested action"}
	error, _ := json.Marshal(unauthorized)

	enableCors(&w)
	// Enabled PREFLIGHT calls
	if r.Method == "OPTIONS" {
		return
	}

	// Get the Jwt
	jwt := r.Header.Get("Authorization")
	query := r.URL.Query()
	t := query["token"]
	bt := query["bearer_token"]

	if jwt == "" {
		if len(t) < 1 || t[0] == "" {
			if len(bt) < 1 || bt[0] == "" {
				errorLogger.Error("Unable to retrieve JWToken from Authorization header or query parameter")
				http.Error(w, string(error), 401)
				return
			}
			jwt = bt[0]
		} else {
			jwt = t[0]
		}
	}

	// Decode it
	mapClaims := make(map[string]interface{})
	jwt = strings.Replace(jwt, "Bearer ", "", 1)
	decoded, jwkset, err := token.Decode(jwt, server.JwkSet, JwtIssuer)
	server.JwkSet = jwkset
	if err != nil {
		errorLogger.Error(err.Error())
		http.Error(w, string(error), 401)
		return
	}
	claims, err := json.Marshal(decoded)
	if err := json.Unmarshal(claims, &mapClaims); err != nil {
		errorLogger.Error(err.Error())
		http.Error(w, string(error), 401)
		return
	}

	if JwtCheckExp {
		// Make sure the exp is before today...
		if _, ok := mapClaims["exp"]; ok != true {
			errorLogger.Error(err.Error())
			http.Error(w, string(error), 401)
			return
		}
		exp := time.Unix(int64(mapClaims["exp"].(float64)), 0)
		now := time.Now()

		if exp.Before(now) {
			errorLogger.Error("Token is expired")
			http.Error(w, string(error), 401)
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

// NewServer creates a new Server object with the jwkset retrieved from the issuer
func NewServer(issuer string) *Server {
	jwks, err := token.JwkSetGet(issuer)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.WithFields(log.Fields{
			"keyset": jwks,
			"issuer": issuer,
		}).Fatal("Unable to retrieve keyset")
	}

	return &Server{
		Issuer: issuer,
		JwkSet: jwks,
	}
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "authorization")
	(*w).Header().Set("Access-Control-Max-Age", "1728000")
	(*w).Header().Set("Access-Control-Expose-Headers", "")
}
