package main

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"../../pkg/httpserver"
	"../../pkg/token"
)

var (
	ListenPortStr     string
	ListenPort        int
	JwtIssuer         string
	JwtOutboundHeader string
	CheckExp          bool
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	ListenPortStr = os.Getenv("LISTEN_PORT")
	var err error
	ListenPort, err = strconv.Atoi(ListenPortStr)
	if err != nil {
		log.Warn("Unable to convert LISTEN_PORT to integer: using 3000")
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

	token.JwtIssuer = JwtIssuer

	// Optional envs
	if JwtOutboundHeader != "" {
		token.JwtOutboundHeader = JwtOutboundHeader
	}

	token.JwtCheckExp = CheckExp
}

func main() {
	log.Info("Starting auth-jwt service")
	server := httpserver.NewServer(JwtIssuer)
	log.Fatal(server.Start(ListenPort))
}
