package main

import (
	"GOVSSO-Mock/app/idtoken"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"time"
)

func (this *routeHandler) handleAuthTokenGeneration(c *gin.Context) {
	var request tokenRequest
	if err := c.ShouldBind(&request); err != nil {
		log.Error().Err(err).Msg("Invalid token request.")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid_request",
			"status_code": "400",
		})
		return
	}

	var authParams *authParams
	if request.Code != "" {
		// Authentication
		authParams = this.authParamsStore.getAndDeleteParams(request.Code)
	} else if request.RefreshToken != "" {
		// Refresh existing session
		authParams = this.authParamsStore.getAndDeleteParams(request.RefreshToken)
	}
	if authParams == nil {
		log.Error().Msg("Invalid code or refresh_token")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid_request",
			"status_code": "400",
		})
		return
	}
	if time.Now().After(authParams.expires) {
		log.Error().Msg("Expired token")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid_request",
			"status_code": "400",
		})
		return
	}
	newRefreshToken := this.generateRandomString()
	this.authParamsStore.addParams(newRefreshToken, *authParams)
	idTokenClaims := this.initAuthIdTokenClaims(*authParams)
	idToken, err := this.idTokenService.CreateAndSignAuthIdToken(idTokenClaims)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create and sign identity token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "server_error",
			"error_debug": "Failed to create and sign identity token",
			"status_code": "500",
		})
		return
	}

	log.Info().Msgf("Generated identity token: %s; generated refresh token: %s", idToken, newRefreshToken)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  "not-used-in-govsso-mock",
		"id_token":      idToken,
		"refresh_token": newRefreshToken,
		"token_type":    "bearer",
		"expires_in":    1,
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

func (this *routeHandler) generateRandomString() string {
	const characterCount = 87
	const characterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-."
	tokenBytes := make([]byte, characterCount)
	for i := range tokenBytes {
		tokenBytes[i] = characterSet[rand.Intn(len(characterSet))]
	}
	return string(tokenBytes)
}

func (this *authParamsStore) getNewTokenExpirationTime() time.Time {
	return time.Now().Add(15 * time.Minute)
}
