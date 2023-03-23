package idtoken

import (
	"github.com/go-jose/go-jose/v3/jwt"
)

type AuthClaims struct {
	jwt.Claims
	Acr                 string   `json:"acr"`
	Amr                 []string `json:"amr"`
	AtHash              string   `json:"at_hash"`
	Birthdate           string   `json:"birthdate"`
	FamilyName          string   `json:"family_name"`
	GivenName           string   `json:"given_name"`
	Nonce               string   `json:"nonce"`
	SessionId           string   `json:"sid"`
	PhoneNumber         string   `json:"phone_number,omitempty"`
	PhoneNumberVerified string   `json:"phone_number_verified,omitempty"`
}
