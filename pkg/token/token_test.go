package token

import (
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

func TestDecode(t *testing.T) {

	var rsaPrivateKey *rsa.PrivateKey
	var signeropts *jose.SignerOptions
	keyid := time.Now().String()
	header := jose.HeaderKey("kid")
	signeropts.WithHeader(header, keyid)
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: rsaPrivateKey}, (nil)).WithType("JWT")
	// sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: rsaPrivateKey}, (signeropts.WithHeader(header, keyid)).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	cl := jwt.Claims{
		Subject:  "admin@example.com",
		Issuer:   "https://localhost/api/v1/oauth2/token",
		Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Audience: jwt.Audience{"leela", "fry"},
	}
	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		panic(err)
	}
	fmt.Println(raw)
	key := jose.JSONWebKey{
		Algorithm: "RS256",
		Key:       rsaPrivateKey.PublicKey,
		KeyID:     keyid,
	}
	keys := []jose.JSONWebKey{key}

	keyset := jose.JSONWebKeySet{Keys: keys}
	// jwt := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IjIwMTgtMDYtMDZUMjI6MzU6MTQrMDA6MDAifQ.eyJ1c2VyX2lkIjoiNjY0Iiwic2NvcGUiOiJhZG1pbi11c2VyOmFsbCIsImV4cGlyZXNfYXQiOiIyMDE4LTA3LTIyVDE5OjM5OjAwWiIsInV1aWQiOiJiMTcwYTU5Ny0xMDE1LTQwNjAtOGY1OS0xMzdjNDYzNzNkNzAiLCJvcmdhbml6YXRpb25faWQiOiI1OCIsImV4cCI6MTUzMjI4ODM0MCwic3ViIjoidGdhbmVtQGFzcGVyYXNvZnQuY29tIiwibmFtZSI6IlRvbSBHYW5lbSIsImdpdmVuX25hbWUiOiJUb20iLCJmYW1pbHlfbmFtZSI6IkdhbmVtIiwiaWF0IjoxNTMyMjQ1MTQwLCJpc3MiOiJodHRwczovL2FwaS5pYm1hc3BlcmEuY29tL2FwaS92MS9vYXV0aDIvdG9rZW4iLCJpYm1pZF9pZCI6IklCTWlkLTUwQzRUV05QSDEiLCJpZCI6ImFvYy02NjQiLCJyZWFsbWlkIjoiYW9jLWlibWlkIiwiaWRlbnRpZmllciI6IjY2NCJ9.l49rDjfBp4PKzE176zbkus7_BeizRvsOGw_N85vHgGHLSiHwK07fX81vJbK1CR6vXU57CsrdXqy6cGnnk5mtLzghQY_meDAxB745t5c3Rx5wbJLe11BtiWE3LN4EvAoqCVdL8BP9gIzJE-n3_nrLcgXVJDKBBvpI0H52JHLvI0vJRHfdAaRzm-GKMUeL1q0VSzFLhf30O-hBWi3Edz_SBrbMlQY2LYUaJ_cMXNmvt_YflXnx5md7mYl5BmCbWeSCsnG2aucDB7tUK3_YQuWDglqgzQrKP2grdOdldDeUSIUKsZU_rJIKIEKRiy6hiBGA6zdMCcTGakZj2cPcR4OkCw"

	// keyjson := `{"keys":[{"kid":"2018-05-23T18:40:16+00:00","alg":"RS256","e":"AQAB","n":"1Tg29euEMj8_dQW3EgYdlhOtDQ277jTVAGcMpn2wXWmHgeWX555T7Ni_N3gPDaRlERaWC7KrZ6oCDV20llZK2B9_uB0VObsnva-39SyRHty4R0iizKHiog0a1htkPsPT7gcZNMmZoWbCdS2T8HQBLUwCRWi0FMsxoCc_M3b9IZJxOU_mwVTuH-wyUovTziLJ8oMaiWoDJVg-veb55ybQyI1unrA4mfaNvK1eXfeXfu8BNpsXh6knpllfw9roCLsVgj4qZVgDyw7MUILn4YTw-sYDsUmjo1DdH34DTkD3nvpTP7BqDZW1w02eZCwsma6rg6FRlPDx6_Rorr4IP22O9Q","kty":"RSA"},{"kid":"2018-06-06T22:35:14+00:00","alg":"RS256","e":"AQAB","n":"3erR0oihAaiS2PXHqFH0vh_fsu_GLBYt05-GsGpSdpkJpTq_7ERH-DdyWqab5Dm_aiP-73TSlYG46vYZMG6AN6y53gWwA5T7cx_FZW950c-VySpG-nwN1mTafM4yefofsQJ4Ues3KP-YaSdeBQqtlrem7x9qWTvoqV-WLiAsp1I8NZEjvj68CxjO32ocDc3o3YtBnZlMMMJeWbxcbD9E6oQNfpiHnzPd2L7IiHfFpf5vcWCnp1epTPRsw9HTBCeQ0mqmhzYoFqGwOxS4F5IPOaS92rHhrmo8nYzH2sEsO3gSJiPdR1jplK9UTei3HFVq4kABNGoHmyJaj5R8KHgnCw","kty":"RSA"}]}`
	// json.Unmarshal([]byte(keyjson), &keyset)
	issuer := "http://localhost/.well-known/jwks.json"
	a, b, err := Decode(raw, keyset, issuer)
	fmt.Println(a)
	fmt.Println(b)
	assert.Nil(t, err)
}

func TestJwkSetGet(t *testing.T) {
}
