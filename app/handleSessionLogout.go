package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
)

func (this *routeHandler) handleSessionLogout(c *gin.Context) {
	var request logoutRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		handleError(c, request, err, "Invalid logout request", http.StatusBadRequest)
		return
	}

	idTokenClaims, err := this.idTokenService.ParseClaimsFromAuthIdToken(request.IdTokenHint)
	if err != nil {
		handleError(c, request, err, "Failed to parse claims from auth identity token", http.StatusBadRequest)
		return
	}

	clientId := idTokenClaims.Audience[0]
	client := this.findFromPredefinedClients(clientId)
	if client == nil {
		handleError(c, request, err, "Invalid OIDC client", http.StatusBadRequest)
		return
	}

	if !client.isValidPostLogoutRedirectUri(request.PostLogoutRedirectUri) {
		handleError(c, request, err, "Invalid post logout redirect URI", http.StatusBadRequest)
		return
	}

	authParams := this.authParamsStore.getAndDeleteBySessionId(idTokenClaims.SessionId)
	var logoutToken string
	if authParams != nil && authParams.isCustomLogoutToken() {
		logoutToken, err = this.idTokenService.CreateAndSignCustomLogoutToken(authParams.LogoutTokenClaims)
		log.Debug().Msgf("Generated logout token: %s; with predefined claims", logoutToken)
	} else {
		logoutToken, err = this.idTokenService.CreateAndSignLogoutToken(clientId, idTokenClaims.SessionId)
		log.Debug().Msgf("Generated logout token: %s", logoutToken)
	}
	if err != nil {
		handleError(c, request, err, "Failed to create and sign logout token", http.StatusInternalServerError)
		return
	}

	go this.backchannelLogout(client.BackchannelLogoutUri, logoutToken, idTokenClaims.SessionId)

	postLogoutRedirectUri := request.PostLogoutRedirectUri
	if request.State != "" {
		postLogoutRedirectUri += "?state=" + url.QueryEscape(request.State)
	}
	c.Redirect(http.StatusMovedPermanently, postLogoutRedirectUri)
}

func handleError(c *gin.Context, request logoutRequest, err error, message string, status int) {
	incidentNr := uuid.New()
	log.Error().Err(err).Msgf("Invalid logout request: %s Incident nr: %s", request, incidentNr)
	c.HTML(status, "error.html", gin.H{
		"error":       "Logout error.",
		"message":     message,
		"incident_nr": incidentNr,
	})
}
