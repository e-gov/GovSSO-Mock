package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
)

func (this *routeHandler) handleBackchannelSessionLogout(c *gin.Context) {
	clientId := c.PostForm("client_id")
	sessionId := c.PostForm("session_id")
	logoutToken, err := this.idTokenService.CreateAndSignLogoutToken(clientId, sessionId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create and sign logout token")
		return
	}

	backchannelLogoutUri := c.PostForm("backchannel_logout_uri")
	backchannelLogoutResult := this.backchannelLogout(backchannelLogoutUri, logoutToken, sessionId)
	c.Redirect(http.StatusMovedPermanently,
		fmt.Sprintf("%s?logout-result-message=%s", this.config.HostUri(), backchannelLogoutResult))
}

func (this *routeHandler) backchannelLogout(backchannelLogoutUri string, logoutToken string, sessionId string) string {
	var resultMessage string
	resp, err := this.httpClient.PostForm(backchannelLogoutUri, url.Values{"logout_token": {logoutToken}})
	if err != nil {
		resultMessage = fmt.Sprintf("Failed to perform back-channel logout request: %s", err)
		log.Error().Err(err).Msg(resultMessage)
		return resultMessage
	}

	if resp.StatusCode != http.StatusOK {
		resultMessage = fmt.Sprintf("Back-channel logout request completed, client application responded with unsuccessful HTTP status code %d.", resp.StatusCode)
		log.Error().Msgf(resultMessage)
		return resultMessage
	}

	resultMessage = fmt.Sprintf("Back-channel logout request completed successfully for session ID '%s'.", sessionId)
	log.Info().Msg(resultMessage)
	return resultMessage
}
