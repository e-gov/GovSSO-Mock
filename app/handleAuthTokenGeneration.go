package main

import (
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

	authParams := this.authParamsStore.getAndDeleteByTokenId(request.getId())
	if authParams == nil {
		log.Error().Msg("Invalid code or refresh_token")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid_request",
			"status_code": "400",
		})
		return
	}
	if authParams.isExpired() {
		log.Error().Msg("Expired token")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid_request",
			"status_code": "400",
		})
		return
	}
	if request.isRefreshTokenRequest() {
		updateClaims(authParams)
	}

	newRefreshToken := this.generateRandomString()
	this.authParamsStore.addParams(newRefreshToken, *authParams)

	idToken, err := this.idTokenService.CreateAndSignAuthIdToken(authParams.IdTokenClaims)

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

func updateClaims(authParams *authParams) {
	issuedAt := time.Now()
	authParams.IdTokenClaims["iat"] = issuedAt.Unix()
	authParams.IdTokenClaims["exp"] = issuedAt.Add(15 * time.Minute).Unix()
	authParams.IdTokenClaims["jti"] = uuid.New()
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
