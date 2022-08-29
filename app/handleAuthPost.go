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
	})

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?state=%s&code=%s",
		c.PostForm("redirect_uri"),
		c.PostForm("state"),
		code))
}
