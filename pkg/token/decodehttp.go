package token

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	logger            = log.New(os.Stdout, "", log.Ltime|log.LUTC)
	JwtCheckExp       = true
	JwtCookieName     = "jwt"
	JwtSecret         = "secret"
	JwtOutboundHeader = "x-jwt-payload"
)

func DecodeHttpHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, r.UserAgent())

	// Get the JWT
	jwt := r.Header.Get("Authorization")
	if jwt == "" {
		logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "No value in header 'Authorization'. Attempting to check cookie")
		jwtCookie, err := r.Cookie(JwtCookieName)
		if err != nil {
			logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "Unable to retrieve jwt from header 'Authorization' or cookie "+JwtCookieName)
			http.Error(w, fmt.Sprintf("Unable to retrieve jwt from header '%s' or cookie '%s'", "Authorization", JwtCookieName), 401)
			return
		}
		jwt = jwtCookie.String()
	}

	// Decode it
	jwt = strings.Replace(jwt, "Bearer ", "", 1)

	mapClaims, err := Decode(jwt, JwtSecret)
	if err != nil {
		logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, err.Error())
		http.Error(w, err.Error(), 401)
		return
	}

	if JwtCheckExp {
		// Make sure the exp is before today...
		if _, ok := mapClaims["exp"]; ok != true {
			logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "exp was not provided in the json payload")
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
			logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "Token is expired")
			http.Error(w, "This token is expired", 401)
			return
		}
	}

	// Put it in the header
	payload, err := json.Marshal(mapClaims)
	if err != nil {
		logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, err.Error())
		http.Error(w, err.Error(), 401)
		return
	}
	logger.Printf("[%s] [%s] [%s %s] Adding payload: \n%s\nTo header: %s", r.RemoteAddr, r.Host, r.Method, r.RequestURI, string(payload), JwtOutboundHeader)
	w.Header().Set(JwtOutboundHeader, string(payload))
	logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "Successfully authorized")
}
