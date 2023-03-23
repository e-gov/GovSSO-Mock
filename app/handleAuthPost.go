package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
)

func (this *routeHandler) handleAuthPost(c *gin.Context) {
	var request authParams
	if err := c.ShouldBind(&request); err != nil {
		log.Error().Err(err).Msg("Invalid auth request.")
		c.JSON(http.StatusOK, gin.H{
			"error":       "invalid_request",
			"status_code": "400",
		})
		return
	}

	code := this.generateRandomString()
	this.authParamsStore.addParams(code, request)
	redirectUri := c.PostForm("redirect_uri")
	if len(redirectUri) == 0 {
		redirectUri = this.findFromPredefinedClients(c.PostForm("client_id")).RedirectUris[0]
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?state=%s&code=%s",
		redirectUri,
		url.QueryEscape(c.PostForm("state")),
		code))
}
