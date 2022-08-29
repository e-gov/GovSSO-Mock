package idtoken

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	idTokenSigningKey *rsa.PrivateKey
	idTokenIssuerUrl  string
	idTokenSignKeyId  string
}

func NewService(idTokenSigningKey *rsa.PrivateKey, idTokenIssuerUrl string, idTokenSignKeyId string) Service {
	return Service{
		idTokenSigningKey: idTokenSigningKey,
		idTokenIssuerUrl:  idTokenIssuerUrl,
		idTokenSignKeyId:  idTokenSignKeyId,
	}
}

func (this Service) CreateAndSignLogoutToken(clientId string, sessionId string) (string, error) {
	logoutToken := jwt.NewWithClaims(jwt.SigningMethodRS256, LogoutClaims{
		ClientId:  []string{clientId},
		IssuedAt:  time.Now().Unix(),
		Issuer:    this.idTokenIssuerUrl,
		Jti:       uuid.New(),
		SessionId: sessionId,
	})
	return logoutToken.SignedString(this.idTokenSigningKey)
}

func (this Service) CreateAndSignAuthIdToken(authClaims AuthClaims) (string, error) {
	idToken := jwt.NewWithClaims(jwt.SigningMethodRS256, authClaims)
	idToken.Header["kid"] = this.idTokenSignKeyId
	return idToken.SignedString(this.idTokenSigningKey)
}

func (this Service) ParseClaimsFromAuthIdToken(idTokenHint string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	if _, err := jwt.ParseWithClaims(idTokenHint, claims, func(*jwt.Token) (interface{}, error) {
		return &this.idTokenSigningKey.PublicKey, nil
	}); err != nil {
		return nil, err
	}
	return claims, nil
}
