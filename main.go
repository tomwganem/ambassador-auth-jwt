package main

import (
	"os"
	"strconv"
	"time"

	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/tomwganem/ambassador-auth-jwt/pkg/httpserver"
)

var (
	// Version should correspond to a git tag
	Version = "0.3.1-rc.0"
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
)

func init() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	raven.SetEnvironment(os.Getenv("SENTRY_CURRENT_ENV"))
	raven.SetRelease(Version)
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
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
	checkExp := os.Getenv("CHECK_EXP")
	allowBasicAuthPassThrough := os.Getenv("ALLOW_BASIC_AUTH_PASSTHROUGH")

	if JwtIssuer == "" {
		log.Fatal("JWT_ISSUER is empty")
	}

	CheckExp = false
	if checkExp != "" {
		b, err := strconv.ParseBool(checkExp)
		if err != nil {
			log.Warn("Unable to convert CHECK_EXP to bool: setting to false")
			b = false
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
	// Optional envs
	if JwtOutboundHeader != "" {
		httpserver.JwtOutboundHeader = JwtOutboundHeader
	}
}

func main() {
	server := httpserver.NewServer(JwtIssuer)
	log.Fatal(server.Start(ListenPort))
}
