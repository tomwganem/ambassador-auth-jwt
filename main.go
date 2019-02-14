package main

import (
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tomwganem/ambassador-auth-jwt/pkg/httpserver"
)

var (
	// Version should correspond to a git tag
	Version = "0.1.0-rc.1"
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
)

func init() {
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.WithFields(log.Fields{
		"version": Version,
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

	httpserver.JwtIssuer = JwtIssuer

	// Optional envs
	if JwtOutboundHeader != "" {
		httpserver.JwtOutboundHeader = JwtOutboundHeader
	}

	httpserver.JwtCheckExp = CheckExp
}

func main() {
	server := httpserver.NewServer(JwtIssuer)
	log.Fatal(server.Start(ListenPort))
}
