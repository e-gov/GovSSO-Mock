package main

import (
	"GOVSSO-Mock/app/idtoken"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (this *routeHandler) handleAuthTokenGeneration(c *gin.Context) {
	code := c.PostForm("code")
	authParams := this.authParamsStore.getParams(code)
	idTokenClaims := this.initAuthIdTokenClaims(*authParams)
	idToken, err := this.idTokenService.CreateAndSignAuthIdToken(idTokenClaims)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create and sign identity token")
		return
	}

	log.Info().Msgf("Generated identity token: %s", idToken)
	c.JSON(http.StatusOK, gin.H{
		"access_token": "not-used-in-govsso-mock",
		"id_token":     idToken,
		"token_type":   "Bearer",
		"expires_in":   3600,
	})
}

func (this *routeHandler) initAuthIdTokenClaims(authParams authParams) idtoken.AuthClaims {
	claims := idtoken.AuthClaims{
		Acr:        authParams.acr,
		Amr:        []string{authParams.amr},
		AtHash:     "QeS-bzYvCt5cN0c0BLVrAA",
		ClientId:   authParams.clientId,
		Birthdate:  authParams.birthdate,
		ExpiresAt:  time.Now().Add(15 * time.Minute).Unix(),
		FamilyName: authParams.familyName,
		GivenName:  authParams.givenName,
		IssuedAt:   time.Now().Unix(),
		Issuer:     this.config.HostUri(),
		Jti:        uuid.New(),
		Nonce:      authParams.nonce,
		Subject:    authParams.subject,
	}

	if len(authParams.phone) != 0 {
		claims.PhoneNumber = authParams.phone
		claims.PhoneNumberVerified = "true"
	}

	if authParams.sessionId == nil {
		claims.SessionId = uuid.New().String()
	} else {
		claims.SessionId = *authParams.sessionId
	}
	return claims
}
