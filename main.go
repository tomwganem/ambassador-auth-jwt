package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/tomwganem/ambassador-auth-jwt/pkg/httpserver"
)

var (
	// Version should correspond to a git tag
	Version = "0.4.0-rc.0"
	// ListenPortStr saves the value extracted from the LISTEN_PORT env var
	ListenPortStr string
	// ListenPort saves LISTEN_PORT as an integer
	ListenPort int
	// JwtIssuer is set by the JWT_ISSUER env variable. It saves the url where the JWKeyset is found
	JwtIssuer string
	// JwtOutboundHeader defaults to X-JWT-PAYLOAD and is returned in the response
	JwtOutboundHeader string
	// CheckExp is a simple flag to check whether tokens are expired
	CheckExp bool
	// AllowBasicAuthPassThrough control whether basic auth requests get rejected or not
	AllowBasicAuthPassThrough bool
	// AllowBasicAuthHeader specifies the header that the basic auth creds are in
	AllowBasicAuthHeaders string
	// AllowBasicAuthPathRegex specifies the path that basic auth requests are allowed on
	AllowBasicAuthPathRegex string
	// NewErrorMessageRegex specifies the paths that we return the new error structure for (needed for backwards compatibility)
	NewErrorMessageRegex string
)

func init() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	raven.SetEnvironment(os.Getenv("SENTRY_CURRENT_ENV"))
	switch logLevel := os.Getenv("LOG_LEVEL"); logLevel {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	raven.SetRelease(Version)
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	log.SetOutput(os.Stdout)
	log.WithFields(log.Fields{
		"version":            Version,
		"sentry_dsn":         os.Getenv("SENTRY_DSN"),
		"sentry_environment": os.Getenv("SENTRY_CURRENT_ENV"),
	}).Info("Starting ambassador-auth-jwt")
	ListenPortStr = os.Getenv("LISTEN_PORT")
	var err error
	ListenPort, err = strconv.Atoi(ListenPortStr)
	if err != nil {
		log.Warn("Unable to convert LISTEN_PORT to integer, defaulting to port 3000")
		ListenPort = 3000
	}

	JwtIssuer = os.Getenv("JWT_ISSUER")
	JwtOutboundHeader = os.Getenv("JWT_OUTBOUND_HEADER")
	AllowBasicAuthHeaders := os.Getenv("ALLOW_BASIC_AUTH_HEADERS")
	AllowBasicAuthPathRegex := os.Getenv("ALLOW_BASIC_AUTH_PATH_REGEX")
	NewErrorMessageRegex := os.Getenv("NEW_ERROR_MESSAGE_REGEX")

	checkExp := os.Getenv("CHECK_EXP")
	allowBasicAuthPassThrough := os.Getenv("ALLOW_BASIC_AUTH_PASSTHROUGH")

	if JwtIssuer == "" {
		log.Fatal("JWT_ISSUER is empty")
	}

	CheckExp = true
	if checkExp != "" {
		b, err := strconv.ParseBool(checkExp)
		if err != nil {
			log.Warn("Unable to convert CHECK_EXP to bool: setting to true")
			b = true
		}
		CheckExp = b
	}

	AllowBasicAuthPassThrough = false
	if allowBasicAuthPassThrough != "" {
		b, err := strconv.ParseBool(allowBasicAuthPassThrough)
		if err != nil {
			log.Warn("Unable to convert ALLOW_BASIC_AUTH_PASSTHROUGH to bool: setting to false")
			b = false
		}
		AllowBasicAuthPassThrough = b
	}

	httpserver.JwtIssuer = JwtIssuer
	httpserver.JwtCheckExp = CheckExp
	httpserver.AllowBasicAuthPassThrough = AllowBasicAuthPassThrough
	if AllowBasicAuthHeaders != "" {
		httpserver.AllowBasicAuthHeaders = strings.Split(AllowBasicAuthHeaders, ",")
	}
	if AllowBasicAuthPathRegex != "" {
		httpserver.AllowBasicAuthPathRegex = regexp.MustCompile(AllowBasicAuthPathRegex)
	}
	if NewErrorMessageRegex != "" {
		httpserver.NewErrorMessageRegex = regexp.MustCompile(NewErrorMessageRegex)
	}
	if JwtOutboundHeader != "" {
		httpserver.JwtOutboundHeader = JwtOutboundHeader
	}
}

func main() {
	server := httpserver.NewServer(JwtIssuer)
	log.Fatal(server.Start(ListenPort))
}
