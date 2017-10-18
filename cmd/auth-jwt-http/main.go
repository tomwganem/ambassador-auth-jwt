package main

import (
	"log"
	"os"
	"strconv"

	"github.com/kminehart/ambassador-auth-jwt/pkg/httpserver"
	"github.com/kminehart/ambassador-auth-jwt/pkg/token"
)

var (
	ListenPortStr     string
	ListenPort        int
	JwtSecret         string
	JwtCookieName     string
	JwtOutboundHeader string
	CheckExp          bool
)

func init() {
	ListenPortStr = os.Getenv("LISTEN_PORT")
	var err error
	ListenPort, err = strconv.Atoi(ListenPortStr)
	if err != nil {
		log.Printf("Error converting LISTEN_PORT to int, using 3000. Error: %s\n", err.Error())
		ListenPort = 3000
	}

	JwtSecret = os.Getenv("JWT_SECRET")
	JwtCookieName = os.Getenv("JWT_SECRET")
	JwtOutboundHeader = os.Getenv("JWT_OUTBOUND_HEADER")
	checkExp := os.Getenv("CHECK_EXP")

	// For security reasons, we will not set a default JWT_SECRET.
	// Generate a random secret and don't share it.
	if JwtSecret == "" {
		log.Fatal("JWT_SECRET is empty")
	}

	CheckExp = false
	if checkExp != "" {
		b, err := strconv.ParseBool(checkExp)
		if err != nil {
			log.Println("Could not convert CHECK_EXP to bool. Defaulting to false Error: %s\n", err.Error())
			b = false
		}
		CheckExp = b
	}

	token.JwtSecret = JwtSecret
	token.JwtCookieName = JwtCookieName
	token.JwtOutboundHeader = JwtOutboundHeader
	token.JwtCheckExp = CheckExp
}

func main() {
	log.Println("Starting auth-jwt service...")
	server := httpserver.NewServer(JwtSecret)
	log.Fatal(server.Start(ListenPort))
}
