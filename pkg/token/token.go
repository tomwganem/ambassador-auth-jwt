package token

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	raven "github.com/getsentry/raven-go"
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
	}).Info("Retrieving Keyset")
	return keyset, nil
}

// Initialize issuer with its jwk set
func JwkSetGetMap(issuers []string) (map[string]jose.JSONWebKeySet, error) {
	keysetIssuerMap := make(map[string]jose.JSONWebKeySet)
	for _, issuer := range issuers {
		keyset := jose.JSONWebKeySet{}
		resp, err := http.Get(issuer)
		if err != nil {
			return keysetIssuerMap, err
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&keyset); err != nil && err != io.EOF {
			return keysetIssuerMap, err
		}
		log.WithFields(log.Fields{
			"keyset": keyset,
			"issuer": issuer,
		}).Info("Retrieving Keyset")
		keysetIssuerMap[issuer] = keyset
	}
	return keysetIssuerMap, nil
}

// Decode the raw token and validate it with a JWK Set.
func Decode(jwtoken string, jwkset jose.JSONWebKeySet, issuer string) (map[string]interface{}, jose.JSONWebKeySet, error) {
	claims := struct {
		*jwt.Claims
		ExpiresAt      string `json:"expires_at,omitempty"`
		Scope          string `json:"scope,omitempty"`
		OrganizationID string `json:"organization_id,omitempty"`
		UUID           string `json:"uuid,omitempty"`
	}{}
	mapClaims := make(map[string]interface{})
	token, err := jwt.ParseSigned(jwtoken)
	if err != nil {
		return mapClaims, jwkset, fmt.Errorf("Could not read jwt")
	}
	keyid := token.Headers[0].KeyID
	jwk := jwkset.Key(keyid)
	if len(jwk) == 0 {
		jwkset, err = JwkSetGet(issuer)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			log.WithFields(log.Fields{
				"issuer": issuer,
				"err":    err,
			}).Fatal("Unable to update keyset")
		}
		log.WithFields(log.Fields{
			"keyset": jwkset,
			"issuer": issuer,
		}).Info("Updating Keyset")
		jwk = jwkset.Key(keyid)
		if len(jwk) == 0 {
			return mapClaims, jwkset, errors.New("Can not find token's key id in jwk set")
		}
	}

	if err := token.Claims(jwk[0].Key.(*rsa.PublicKey), &claims); err != nil {
		raven.CaptureError(err, nil)
		return mapClaims, jwkset, err
	}
	marshalClaims, err := json.Marshal(claims)
	if err != nil {
		raven.CaptureError(err, nil)
		return mapClaims, jwkset, err
	}
	if err := json.Unmarshal(marshalClaims, &mapClaims); err != nil {
		raven.CaptureError(err, nil)
		return mapClaims, jwkset, err
	}

	return mapClaims, jwkset, nil
}
