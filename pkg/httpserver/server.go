package httpserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
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
	// JwtIssuer is the url where we can retrieve a set of public keys to verify rsa based tokens with
	JwtIssuer = "http://localhost/.well-known/jwks.json"
	// JwtOutboundHeader is the name of header the parsed token claims will be inserted into
	JwtOutboundHeader = "X-JWT-PAYLOAD"
	// AllowBasicAuthPassThrough will allow requests with a basic auth authorization header to be passed through
	AllowBasicAuthPassThrough = false
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
	q, _ := url.ParseQuery(r.URL.RawQuery)
	successFields := log.Fields{
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
		"method":      r.Method,
		"path":        r.URL.Path,
		"query":       q,
		"user_agent":  r.UserAgent(),
		"status":      "200",
	}
	errorFields := log.Fields{
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
		"method":      r.Method,
		"path":        r.URL.Path,
		"query":       q,
		"user_agent":  r.UserAgent(),
		"status":      "401",
	}
	successLogger := log.WithFields(successFields)
	errorLogger := log.WithFields(errorFields)

	unauthorized := map[string]string{"code": "unauthorized", "message": "You are not authorized to perform the requested action"}
	error, _ := json.Marshal(unauthorized)

	enableCors(&w)
	// Enabled PREFLIGHT calls
	if r.Method == "OPTIONS" {
		successLogger.Info("CORS Request OK")
		return
	}

	// Get the Jwt
	auth := r.Header.Get("Authorization")
	matched, err := regexp.MatchString(`^Basic*`, auth)

	if matched && AllowBasicAuthPassThrough {
		auth = strings.Replace(auth, "Basic ", "", 1)
		payload, _ := base64.StdEncoding.DecodeString(auth)
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			errorLogger.Error("Authorization Failed")
			http.Error(w, string(error), 401)
			return
		}
		successLogger.Info("Basic Auth Request OK")
		return
	}

	query := r.URL.Query()
	t := query["token"]
	bt := query["bearer_token"]

	if auth == "" {
		if len(t) < 1 || t[0] == "" {
			if len(bt) < 1 || bt[0] == "" {
				errorLogger.Error("Unable to retrieve JWToken from Authorization header or query parameter")
				http.Error(w, string(error), 401)
				return
			}
			auth = bt[0]
		} else {
			auth = t[0]
		}
	}

	claims := make(map[string]interface{})
	auth = strings.Replace(auth, "Bearer ", "", 1)
	claims, jwkset, err := token.Decode(auth, server.JwkSet, JwtIssuer)
	server.JwkSet = jwkset
	if err != nil {
		errorLogger.Error(err.Error())
		http.Error(w, string(error), 401)
		return
	}
	exp := time.Now()
	if JwtCheckExp {
		// Make sure the exp is before today...
		if _, ok := claims["exp"]; ok != true {
			if _, ok := claims["expires_at"]; ok != true {
				errorLogger.Error(err.Error())
				http.Error(w, string(error), 401)
				return
			} else {
				exp, err = time.Parse(time.RFC3339, claims["expires_at"].(string))
				if err != nil {
					raven.CaptureError(err, nil)
					errorLogger.Error(err.Error())
					http.Error(w, string(error), 500)
				}
			}
		} else {
			exp = time.Unix(int64(claims["exp"].(float64)), 0)
		}

		now := time.Now()

		if exp.Before(now) {
			errorLogger.Error("Token is expired")
			http.Error(w, string(error), 401)
			return
		}
	}
	marshaledClaims, err := json.Marshal(claims)
	successFields["claims"] = claims
	log.WithFields(successFields).Info("Authentication Success")
	w.Header().Set(JwtOutboundHeader, string(marshaledClaims))
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
