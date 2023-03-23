package idtoken

import (
	"crypto/rsa"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	idTokenSigningKey *jose.SigningKey
	idTokenIssuerUrl  string
	idTokenSignKeyId  string
}

func NewService(idTokenSigningKey *jose.SigningKey, idTokenIssuerUrl string, idTokenSignKeyId string) Service {
	return Service{
		idTokenSigningKey: idTokenSigningKey,
		idTokenIssuerUrl:  idTokenIssuerUrl,
		idTokenSignKeyId:  idTokenSignKeyId,
	}
}

func (this Service) CreateAndSignCustomLogoutToken(claims map[string]interface{}) (string, error) {
	var signerOpts = jose.SignerOptions{}
	signerOpts.WithHeader("kid", this.idTokenSignKeyId)
	signerOpts.WithHeader("typ", "JWT")
	signer, err := jose.NewSigner(*this.idTokenSigningKey, &signerOpts)
	if err != nil {
		return "", err
	}

	return jwt.Signed(signer).Claims(claims).CompactSerialize()
}

func (this Service) CreateAndSignLogoutToken(clientId string, sessionId string) (string, error) {
	claims := map[string]interface{}{
		"aud": []string{clientId},
		"events": map[string]interface{}{
			"http://schemas.openid.net/event/backchannel-logout": struct{}{},
		},
		"iat": time.Now().Unix(),
		"iss": this.idTokenIssuerUrl,
		"jti": uuid.New(),
		"sid": sessionId,
	}

	return this.CreateAndSignCustomLogoutToken(claims)
}

func (this Service) CreateAndSignAuthIdToken(claims map[string]interface{}) (string, error) {
	var signerOpts = jose.SignerOptions{}
	signerOpts.WithHeader("kid", this.idTokenSignKeyId)
	signerOpts.WithHeader("typ", "JWT")
	signer, err := jose.NewSigner(*this.idTokenSigningKey, &signerOpts)
	if err != nil {
		return "", err
	}
	return jwt.Signed(signer).Claims(claims).CompactSerialize()
}

func (this Service) ParseClaimsFromAuthIdToken(idTokenHint string) (*AuthClaims, error) {
	token, err := jwt.ParseSigned(idTokenHint)
	if err != nil {
		return nil, err
	}
	key := this.idTokenSigningKey.Key.(*rsa.PrivateKey)
	claims := &AuthClaims{}
	err = token.Claims(&key.PublicKey, &claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
