package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (this *routeHandler) handleAuthGet(c *gin.Context) {
	if c.Query("prompt") == "none" {
		this.authenticateByIdTokenHintWithoutAuthForm(c)
		return
	}

	subject := c.Query("auto_login")
	if subject != "" {
		this.authenticateBySubjectWithoutAuthForm(c, subject)
	}

	c.HTML(http.StatusOK, "userAuth.html", gin.H{
		"Request": gin.H{
			"ClientId":     c.Query("client_id"),
			"Nonce":        c.Query("nonce"),
			"RedirectUri":  c.Query("redirect_uri"),
			"ResponseType": c.Query("response_type"),
			"Scope":        c.Query("scope"),
			"State":        c.Query("state"),
			"UILocales":    c.Query("ui_locales"),
			"AcrValues":    c.Query("acr_values"),
		},
		"PredefinedUsers": this.predefinedUsers,
		"isPhoneScope":    isPhoneScope(c.Query("scope")),
	})

}

func isPhoneScope(s string) bool {
	for _, v := range strings.Split(s, " ") {
		if v == "phone" {
			return true
		}
	}
	return false
}

func (this *routeHandler) authenticateByIdTokenHintWithoutAuthForm(c *gin.Context) {
	idTokenHint := c.Query("id_token_hint")
	idTokenClaims, err := this.idTokenService.ParseClaimsFromAuthIdToken(idTokenHint)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse claims from auth identity token")
		return
	}

	code := this.authParamsStore.addParams(authParams{
		acr:        idTokenClaims.Acr,
		amr:        idTokenClaims.Amr[0],
		birthdate:  idTokenClaims.Birthdate,
		clientId:   idTokenClaims.ClientId,
		givenName:  idTokenClaims.GivenName,
		familyName: idTokenClaims.FamilyName,
		nonce:      c.Query("nonce"),
		state:      c.Query("state"),
		subject:    idTokenClaims.Subject,
		sessionId:  &idTokenClaims.SessionId,
	})

	this.setCorsHeader(c)
	c.Redirect(http.StatusFound, fmt.Sprintf("%s?state=%s&code=%s",
		c.Query("redirect_uri"),
		c.Query("state"),
		code))
}

func (this *routeHandler) authenticateBySubjectWithoutAuthForm(c *gin.Context, subject string) {
	params := authParams{
		clientId: []string{c.Query("client_id")},
		nonce:    c.Query("nonce"),
		state:    c.Query("state"),
	}

	if user := this.findFromPredefinedUsers(subject); user != nil {
		params.acr = user.Acr
		params.amr = user.Amr
		params.birthdate = user.Birthdate
		params.givenName = user.GivenName
		params.familyName = user.FamilyName
		params.subject = user.Subject
		params.phone = user.Phone
	} else {
		params.acr = "high"
		params.amr = "idcard"
		params.birthdate = "2000-01-01"
		params.givenName = fmt.Sprintf("%s-GivenName", subject)
		params.familyName = fmt.Sprintf("%s-FamilyName", subject)
		params.subject = subject
	}

	code := this.authParamsStore.addParams(params)

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?state=%s&code=%s",
		c.Query("redirect_uri"),
		c.Query("state"),
		code))
}

func (this *routeHandler) findFromPredefinedUsers(subject string) *user {
	for _, user := range this.predefinedUsers {
		if user.Subject == subject {
			return &user
		}
	}
	return nil
}
