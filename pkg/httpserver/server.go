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
	// AllowBasicAuthHeader specifies the header to extract the basic auth request from
	AllowBasicAuthHeader = "Authorization"
	// AllowBasicAuthPathRegex allows us to only allow a basic auth pass through for requests with a certain path
	AllowBasicAuthPathRegex = regexp.MustCompile(`^\/.*`)
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
	debugFields := log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  q,
	}
	successLogger := log.WithFields(successFields)
	errorLogger := log.WithFields(errorFields)
	debugLogger := log.WithFields(debugFields)

	unauthorized := map[string]string{"code": "unauthorized", "message": "You are not authorized to perform the requested action"}
	error, _ := json.Marshal(unauthorized)

	enableCors(&w)
	// Enabled PREFLIGHT calls
	if r.Method == "OPTIONS" {
		successLogger.Info("CORS Request OK")
		return
	}

	auth := r.Header.Get("Authorization")
	query := r.URL.Query()
	t := query["token"]
	bt := query["bearer_token"]
	basicAuthAllowed, msg := basicAuthPassCheck(r, debugLogger)

	// The following checks for tokens passed as a query parameter
	if auth == "" && !basicAuthAllowed {
		if len(t) < 1 || t[0] == "" {
			if len(bt) < 1 || bt[0] == "" {
				errorLogger.Warn("Unable to retrieve JWToken from Authorization header or query parameter. " + msg)
				http.Error(w, string(error), 401)
				return
			}
			auth = bt[0]
		} else {
			auth = t[0]
		}

	} else if auth == "" && basicAuthAllowed {
		log.WithFields(successFields).Info(msg)
		return
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

// basicAuthPassCheck returns a boolean. It will return true if:
// 1. ALLOW_BASIC_AUTH_PASSTHROUGH is set to true
// 2. the path of the request matches ALLOW_BASIC_AUTH_PATH_REGEX
// 3. the value passed in ALLOW_BASIC_AUTH_HEADER looks like a valid basic auth value (i.e. starts with "Basic", can be base64 decoded, can be split into a username:password pair)
func basicAuthPassCheck(r *http.Request, debugLogger *log.Entry) (bool, string) {
	basicAuthRegex := regexp.MustCompile(`^Basic\ *`)
	debugLogger.Debug(fmt.Sprintf("ALLOW_BASIC_AUTH_PASSTHROUGH set to %t", AllowBasicAuthPassThrough))
	if AllowBasicAuthPassThrough {
		debugLogger.Trace(fmt.Sprintf("ALLOW_BASIC_AUTH_PATH_REGEX set to: %s", AllowBasicAuthPathRegex))
		debugLogger.Trace(fmt.Sprintf("ALLOW_BASIC_AUTH_HEADER set to: %s", AllowBasicAuthHeader))
		matchedPath := AllowBasicAuthPathRegex.Match([]byte(r.URL.Path))
		basicAuth := r.Header.Get(AllowBasicAuthHeader)
		matchedAuth := basicAuthRegex.Match([]byte(basicAuth))
		if matchedAuth {
			debugLogger.Trace(fmt.Sprintf("header: %s, does have a value that includes: %s", AllowBasicAuthHeader, basicAuthRegex))
			if matchedPath {
				debugLogger.Trace(fmt.Sprintf("request path: %s, correctly matches regex: %s", r.URL.Path, AllowBasicAuthPathRegex))
				basicAuth = strings.Replace(basicAuth, "Basic ", "", 1)
				payload, err := base64.StdEncoding.DecodeString(basicAuth)
				if err != nil {
					debugLogger.Trace(fmt.Sprintf("basic auth value: %s can not be base64 decoded", basicAuth))
					return false, "Basic Auth Not Allowed"
				}
				pair := strings.SplitN(string(payload), ":", 2)
				if len(pair) == 2 {
					return true, "Basic Auth Allowed"
				}
				debugLogger.Trace(fmt.Sprintf("decoded basic auth value: %s is unable to be split into a username/password pair", payload))
			} else {
				debugLogger.Trace(fmt.Sprintf("request path: %s, does not match regex: %s", r.URL.Path, AllowBasicAuthPathRegex))
			}
		} else {
			debugLogger.Trace(fmt.Sprintf("header: %s, does not have value that matches: %s", AllowBasicAuthHeader, basicAuthRegex))
		}
	}
	return false, "Basic Auth Not Allowed"
}
