package token

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

// JwkSetGet will call the url provided JWT_ISSUER and retreive a JWK Set.
func JwkSetGet(issuer string) ([]jose.JSONWebKey, error) {
	resp, err := http.Get(issuer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Info("Retreiving Keyset from", issuer)
	keyset := jose.JSONWebKeySet{}
	if err := json.NewDecoder(resp.Body).Decode(&keyset); err != nil && err != io.EOF {
		return nil, err
	}
	return keyset.Keys, nil
}

// Decode the raw token and validate it with a JWK Set.
func Decode(jwtoken string, issuer string) (jwt.Claims, error) {
	claims := jwt.Claims{}
	token, err := jwt.ParseSigned(jwtoken)
	if err != nil {
		return claims, fmt.Errorf("Could not read jwt")
	}

	jwks, err := JwkSetGet(issuer)
	for _, jwk := range jwks {
		key := jwk.Key.(*rsa.PublicKey)
		if err := token.Claims(key, &claims); err == nil {
			return claims, nil
		}
	}

	if err != nil {
		return claims, err
	}

	return claims, nil
}
