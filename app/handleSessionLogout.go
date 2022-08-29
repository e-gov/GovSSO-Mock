package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (this *routeHandler) handleSessionLogout(c *gin.Context) {
	idTokenHint := c.Query("id_token_hint")
	idTokenClaims, err := this.idTokenService.ParseClaimsFromAuthIdToken(idTokenHint)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse claims from auth identity token")
		return
	}

	clientId := idTokenClaims.ClientId[0]
	logoutToken, err := this.idTokenService.CreateAndSignLogoutToken(clientId, idTokenClaims.SessionId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create and sign logout token")
		return
	}

	postLogoutRedirectUri := c.Query("post_logout_redirect_uri")
	state := c.Query("state")
	if state != "" {
		postLogoutRedirectUri += "?state=" + state
	}

	// If client is present in predefined clients configuration (./config/clients.json),
	// then also call back-channel logout asynchronously.
	client := this.findFromPredefinedClients(clientId)
	if client != nil {
		go this.backchannelLogout(client.BackchannelLogoutUri, logoutToken, idTokenClaims.SessionId)
	}
	c.Redirect(http.StatusMovedPermanently, postLogoutRedirectUri)
}

func (this *routeHandler) findFromPredefinedClients(clientId string) *client {
	for _, client := range this.predefinedClients {
		if client.ClientId == clientId {
			return &client
		}
	}
	return nil
}
