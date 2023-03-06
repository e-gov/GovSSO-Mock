package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
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

	var request authRequest
	if err := c.ShouldBindQuery(&request); err == nil {
		c.HTML(http.StatusOK, "userAuth.html", gin.H{
			"Request":         request,
			"PredefinedUsers": this.predefinedUsers,
			"isPhoneScope":    strings.Contains(c.Query("scope"), "phone"),
		})
	} else if validationErr, ok := err.(validator.ValidationErrors); ok {
		incidentNr := uuid.New()
		log.Error().Err(err).Msgf("Invalid authorization request. Incident nr: %s", incidentNr)
		handleValidationError(c, validationErr, incidentNr)
	}
}

func handleValidationError(c *gin.Context, verr validator.ValidationErrors, incidentNr uuid.UUID) {
	validationError := verr[0]

	if validationError.Tag() == "valid_response_type" {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s?error=unsupported_response_type&"+
			"error_description=The+authorization+server+does+not+support+obtaining+a+token+using+this+method.+"+
			"The+client+is+not+allowed+to+request+response_type+'%s'&state=%s",
			c.Query("redirect_uri"),
			url.QueryEscape(c.Query("response_type")),
			url.QueryEscape(c.Query("state"))))
	} else if validationError.Tag() == "valid_scope" {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s?error=invalid_scope&The+requested+scope+is+invalid,"+
			"+unknown,+or+malformed.+The+OAuth+2.0+Client+is+not+allowed+to+request+scope+'%s'&state=%s",
			c.Query("redirect_uri"),
			url.QueryEscape(c.Query("scope")),
			url.QueryEscape(c.Query("state"))))
	} else if validationError.Tag() == "valid_state" {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s?error=invalid_state&The+state+is+missing+or+does+not"+
			"+have+enough+characters+and+is+therefore+considered+too+weak.+Request+parameter+'state'+must+be"+
			"+at+least+be+8+characters+long+to+ensure+sufficient+entropy.&state=%s",
			c.Query("redirect_uri"),
			url.QueryEscape(c.Query("state"))))
	} else if validationError.Tag() == "valid_client" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error":       "Authentication error.",
			"message":     "Invalid OIDC client.",
			"incident_nr": incidentNr,
		})
	} else {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error":       "Authentication error.",
			"message":     "Invalid OIDC request.",
			"incident_nr": incidentNr,
		})
	}
}

func (this *routeHandler) authRequestValidation(sl validator.StructLevel) {
	request := sl.Current().Interface().(authRequest)
	authClient := this.findFromPredefinedClients(request.ClientId)

	if authClient == nil {
		sl.ReportError(request.ClientId, "client_id", "ClientId", "valid_client", "")
	} else if !authClient.isValidRedirectUri(request.RedirectUri) {
		sl.ReportError(request.RedirectUri, "redirect_uri", "RedirectUri", "valid_redirect_uri", "")
	}
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
		url.QueryEscape(c.Query("state")),
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
		url.QueryEscape(c.Query("state")),
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
