package token

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	JwtCheckExp       = true
	JwtIssuer         = "secret"
	JwtOutboundHeader = "X-JWT-PAYLOAD"
)

func DecodeHTTPHandler(w http.ResponseWriter, r *http.Request) {
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
		errorLogger.Error("Unable to retrieve Jwt from header 'Authorization'")
		http.Error(w, fmt.Sprintf("Unable to retrieve Jwt from header '%s'", "Authorization"), 401)
		return
	}

	// Decode it
	mapClaims := make(map[string]interface{})
	jwt = strings.Replace(jwt, "Bearer ", "", 1)
	decoded, err := Decode(jwt, JwtIssuer)
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

		exp, err := time.Parse(time.RFC3339, mapClaims["exp"].(string))
		if err != nil {
			logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, err.Error())
			http.Error(w, err.Error(), 401)
			return
		}

		exp = exp.UTC()
		now := time.Now().UTC()

		if exp.Before(now) {
			errorLogger.Error("Token is expired")
			http.Error(w, "This token is expired", 401)
			return
		}
	}
	h := make(map[string]interface{})
	h[JwtOutboundHeader] = mapClaims
	log.WithFields(log.Fields{
		"remote_addr": r.RemoteAddr,
		"host":        r.Host,
		"method":      r.Method,
		"request_uri": r.RequestURI,
		"user_agent":  r.UserAgent(),
		"header":      h,
		"status":      "200",
	}).Info("Authentication Success")
	w.Header().Set(JwtOutboundHeader, string(claims))
}
