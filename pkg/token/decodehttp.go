package token

import (
	"encoding/json"
	"fmt"
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
			return
		}
		jwt = jwtCookie.String()
	}

	// Decode it
	jwt = strings.Replace(jwt, "Bearer ", "", 1)

	mapClaims, err := Decode(jwt, JwtSecret)
	if err != nil {
		logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, err.Error())
		return
	}

	if JwtCheckExp {
		// Make sure the exp is before today...
		if _, ok := mapClaims["exp"]; ok != true {
			logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "exp was not provided in the json payload")
			return
		}

		exp, err := time.Parse(time.RFC3339, mapClaims["exp"].(string))
		if err != nil {
			logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, err.Error())
			return
		}

		exp = exp.UTC()
		now := time.Now().UTC()

		if exp.Before(now) {
			logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "Token is expired")
			return
		}
	}

	// Put it in the header
	payload, err := json.Marshal(mapClaims)
	if err != nil {
		logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, err.Error())
		return
	}
	logger.Printf("[%s] [%s] [%s %s] Adding payload: \n%s\nTo header: %s", r.RemoteAddr, r.Host, r.Method, r.RequestURI, string(payload), JwtOutboundHeader)
	w.Header().Set(JwtOutboundHeader, string(payload))
	logger.Printf("[%s] [%s] [%s %s] %s\n", r.RemoteAddr, r.Host, r.Method, r.RequestURI, "Successfully authorized")
}
