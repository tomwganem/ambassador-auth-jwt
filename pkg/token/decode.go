package token

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

// JwkSetGet will call the url provided JWT_ISSUER and retreive a JWK Set.
func JwkSetGet(issuer string) (jose.JSONWebKeySet, error) {
	keyset := jose.JSONWebKeySet{}
	resp, err := http.Get(issuer)
	if err != nil {
		return keyset, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&keyset); err != nil && err != io.EOF {
		return keyset, err
	}
	log.WithFields(log.Fields{
		"keyset": keyset,
		"issuer": issuer,
	}).Info("Retreiving Keyset")
	return keyset, nil
}

// Decode the raw token and validate it with a JWK Set.
func Decode(jwtoken string, jwkset jose.JSONWebKeySet, issuer string) (jwt.Claims, jose.JSONWebKeySet, error) {
	claims := jwt.Claims{}
	token, err := jwt.ParseSigned(jwtoken)
	if err != nil {
		return claims, jwkset, fmt.Errorf("Could not read jwt")
	}
	keyid := token.Headers[0].KeyID
	jwk := jwkset.Key(keyid)
	if len(jwk) == 0 {
		jwkset, err = JwkSetGet(issuer)
		if err != nil {
			log.WithFields(log.Fields{
				"issuer": issuer,
			}).Fatal("Unable to update keyset")
		}
		log.WithFields(log.Fields{
			"keyset": jwkset,
			"issuer": issuer,
		}).Info("Updating Keyset")
		jwk = jwkset.Key(keyid)
		if len(jwk) == 0 {
			return claims, jwkset, errors.New("Can not find token's key id in jwk set")
		}
	}

	if err := token.Claims(jwk[0].Key.(*rsa.PublicKey), &claims); err != nil {
		return claims, jwkset, err
	}

	return claims, jwkset, nil
}
