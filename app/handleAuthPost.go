package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *routeHandler) handleAuthPost(c *gin.Context) {
	code := this.authParamsStore.addParams(authParams{
		acr:        c.PostForm("acr"),
		amr:        c.PostForm("amr"),
		birthdate:  c.PostForm("birthdate"),
		clientId:   []string{c.PostForm("client_id")},
		givenName:  c.PostForm("given_name"),
		familyName: c.PostForm("family_name"),
		nonce:      c.PostForm("nonce"),
		state:      c.PostForm("state"),
		subject:    c.PostForm("subject"),
		scope:      c.PostForm("scope"),
		phone:      c.PostForm("phone"),
	})

	redirectUri := c.PostForm("redirect_uri")
	if len(redirectUri) == 0 {
		redirectUri = this.findFromPredefinedClients(c.PostForm("client_id")).RedirectUris[0]
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?state=%s&code=%s",
		redirectUri,
		c.PostForm("state"),
		code))
}
