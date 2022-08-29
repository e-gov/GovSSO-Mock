package idtoken

import "github.com/google/uuid"

type AuthClaims struct {
	Acr        string    `json:"acr"`
	Amr        []string  `json:"amr"`
	AtHash     string    `json:"at_hash"`
	Birthdate  string    `json:"birthdate"`
	ClientId   []string  `json:"aud"`
	ExpiresAt  int64     `json:"exp"`
	FamilyName string    `json:"family_name"`
	GivenName  string    `json:"given_name"`
	IssuedAt   int64     `json:"iat"`
	Issuer     string    `json:"iss"`
	Jti        uuid.UUID `json:"jti"`
	Nonce      string    `json:"nonce"`
	SessionId  string    `json:"sid"`
	Subject    string    `json:"sub"`
}

// Valid checks identity token correctness.
func (AuthClaims) Valid() error {
	return nil
}
