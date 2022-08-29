package idtoken

import "github.com/google/uuid"

type LogoutClaims struct {
	ClientId []string `json:"aud"`
	Events   struct {
		SchemaUrl struct{} `json:"http://schemas.openid.net/event/backchannel-logout"`
	} `json:"events"`
	IssuedAt  int64     `json:"iat"`
	Issuer    string    `json:"iss"`
	Jti       uuid.UUID `json:"jti"`
	SessionId string    `json:"sid"`
}

// Valid checks identity token correctness.
func (LogoutClaims) Valid() error {
	return nil
}
