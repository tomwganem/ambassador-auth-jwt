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

// Error is returned in ErrorMsg
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ErrorMsg is returned on unauthorized requests
type ErrorMsg struct {
	StatusCode int     `json:"status_code,omitempty"`
	Errors     []Error `json:"errors"`
}

var (
	// JwtCheckExp will determine if we need to verify if the token is expired or not
	JwtCheckExp = true
	// JwtIssuer is the url where we can retrieve a set of public keys to verify rsa based tokens with
	JwtIssuer = map[string]string{
		"default": "http://localhost/.well-known/jwks.json",
	}
	// JwtOutboundHeader is the name of header the parsed token claims will be inserted into
	JwtOutboundHeader = "X-JWT-PAYLOAD"
	// AllowBasicAuthPassThrough will allow requests with a basic auth authorization header to be passed through
	AllowBasicAuthPassThrough = false
	// AllowBasicAuthHeaders specifies the header to extract the basic auth request from
	AllowBasicAuthHeaders = []string{"Authorization"}
	// AllowBasicAuthPathRegex allows us to only allow a basic auth pass through for requests with a certain path
	AllowBasicAuthPathRegex = regexp.MustCompile(`^\/.*`)
	// NewErrorMessageRegex is used to apply new error structure only to specific endpoints (see bcab0f12220b1e26c6dad10ee6aa1825172051ab). This is for backward compatibility.
	NewErrorMessageRegex = regexp.MustCompile(`^\/.*`)
	// BasicAuthRegex is for checking if a basic auth request is formatted correctly
	BasicAuthRegex = regexp.MustCompile(`^Basic\ .*`)
)

// Server needs to know about the Issuer url to verify tokens against
type Server struct {
	IssuerJwkSetMap map[string]jose.JSONWebKeySet
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

	unauthorizedOld := map[string]string{"code": "unauthorized", "message": "You are not authorized to perform the requested action"}

	unauthorizedNew := ErrorMsg{
		StatusCode: 401,
		Errors: []Error{
			{
				Code:    "unauthorized",
				Message: "You are not authorized to perform the requested action",
			},
		},
	}

	matchedPath := NewErrorMessageRegex.Match([]byte(r.URL.Path))
	var error []byte
	if matchedPath {
		error, _ = json.Marshal(unauthorizedNew)
	} else {
		error, _ = json.Marshal(unauthorizedOld)
	}

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
	matchedAuth := BasicAuthRegex.Match([]byte(auth))
	// Allows basic auth credentials in the Authorization header to be passed through
	if matchedAuth && basicAuthAllowed {
		successLogger.Info(msg)
		return
	}
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
	found,issuer := getJwtIssuer(r)
	if !found {
		errorLogger.Error("Could not find jwt issuer for path " + r.URL.Path)
		http.Error(w, string(error), 401)
		return
	}
	claims, jwkset, err := token.Decode(auth, server.IssuerJwkSetMap[issuer], issuer)
	server.IssuerJwkSetMap[issuer] = jwkset
	if err != nil {
		errorLogger.Error(err.Error())
		http.Error(w, string(error), 401)
		return
	}
	exp := time.Now()
	if JwtCheckExp {
		// Checks to see if the there is an "exp" field
		if _, ok := claims["exp"]; ok != true {
			// Checks to see if there is an "expires_at" field. Note: "expires_at" doesn't follow the RFC and shouldn't be a field in most JWTokens. It's the same as "exp", except it's in RFC3339.
			if _, ok := claims["expires_at"]; ok != true {
				errorLogger.Error(err.Error())
				http.Error(w, string(error), 401)
				return
			}
			exp, err = time.Parse(time.RFC3339, claims["expires_at"].(string))
			if err != nil {
				raven.CaptureError(err, nil)
				errorLogger.Error(err.Error())
				http.Error(w, string(error), 500)
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
func NewServer(issuers []string) *Server {
	jwks, err := token.JwkSetGetMap(issuers)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.WithFields(log.Fields{
			"issuerJwkSetMap": jwks,
			"err":    err,
		}).Fatal("Unable to retrieve keyset")
	}

	return &Server{
		IssuerJwkSetMap: jwks,
	}
}

// enableCors sets some hardcoded headers for OPTIONS requests. We return all OPTIONS requests with a 200.
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
// 3. the value passed in ALLOW_BASIC_AUTH_HEADERS looks like a valid basic auth value (i.e. starts with "Basic", can be base64 decoded, can be split into a username:password pair)
func basicAuthPassCheck(r *http.Request, debugLogger *log.Entry) (bool, string) {
	debugLogger.Trace(fmt.Sprintf("ALLOW_BASIC_AUTH_PASSTHROUGH set to %t", AllowBasicAuthPassThrough))
	if AllowBasicAuthPassThrough {
		debugLogger.Trace(fmt.Sprintf("ALLOW_BASIC_AUTH_HEADERS have values: %v", AllowBasicAuthHeaders))
		for _, header := range AllowBasicAuthHeaders {
			b, msg := basicAuthHeaderCheck(header, r, debugLogger)
			if b {
				return b, msg
			}
		}
	}
	return false, "Basic Auth Not Allowed"
}

func basicAuthHeaderCheck(header string, r *http.Request, debugLogger *log.Entry) (bool, string) {
	debugLogger.Trace(fmt.Sprintf("Checking value in header: %s", header))
	path := r.URL.Path
	debugLogger.Trace(fmt.Sprintf("ALLOW_BASIC_AUTH_PATH_REGEX set to: %s", AllowBasicAuthPathRegex))
	matchedPath := AllowBasicAuthPathRegex.Match([]byte(path))
	basicAuth := r.Header.Get(header)
	matchedAuth := BasicAuthRegex.Match([]byte(basicAuth))
	if matchedAuth {
		debugLogger.Trace(fmt.Sprintf("header: %s, does have a value that includes: %s", header, BasicAuthRegex))
		if matchedPath {
			debugLogger.Trace(fmt.Sprintf("request path: %s, correctly matches regex: %s", path, AllowBasicAuthPathRegex))
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
			debugLogger.Trace(fmt.Sprintf("request path: %s, does not match regex: %s", path, AllowBasicAuthPathRegex))
		}
	} else {
		debugLogger.Trace(fmt.Sprintf("header: %s, does not have value that matches: %s", header, BasicAuthRegex))
	}
	return false, "Basic Auth Not Allowed"
}

func getJwtIssuer(r *http.Request) (bool,string) {
	 path := r.URL.Path
	 for jwt_path, issuer := range JwtIssuer {
		 if (strings.Contains(path, jwt_path)) {
			 return true,issuer
		 }
	 }
	 return false,"Path not found in jwt keys"
}
