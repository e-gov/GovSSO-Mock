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
	"time"
)

func (this *routeHandler) handleAuthGet(c *gin.Context) {
	subject := c.Query("auto_login")
	if subject != "" {
		this.authenticateBySubjectWithoutAuthForm(c, subject)
	}

	var request authRequest
	if err := c.ShouldBindQuery(&request); err == nil {
		c.HTML(http.StatusOK, "userAuth.html", gin.H{
			"Request":         request,
			"Issuer":          this.config.HostUri(),
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

func (this *routeHandler) authenticateBySubjectWithoutAuthForm(c *gin.Context, subject string) {
	claims := map[string]interface{}{
		"at_hash": "YjllNTNjMDY1Y2MxZGNlYmQwODZiZDQwZDkzNzRjNGNjZDQ3YWFlMjgzN2IwZTQ1NTcxODlhMTU4NzhiOWE4Nw==",
		"aud":     []string{c.Query("client_id")},
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     this.config.HostUri(),
		"jti":     uuid.New(),
		"nonce":   c.Query("nonce"),
		"sid":     uuid.New().String(),
	}

	if user := this.findFromPredefinedUsers(subject); user != nil {
		claims["acr"] = user.Acr
		claims["amr"] = []string{user.Amr}
		claims["birthdate"] = user.Birthdate
		claims["given_name"] = user.GivenName
		claims["family_name"] = user.FamilyName
		claims["subject"] = user.Subject
	} else {
		claims["acr"] = "high"
		claims["amr"] = []string{"idcard"}
		claims["birthdate"] = "2000-01-01"
		claims["given_name"] = fmt.Sprintf("%s-GivenName", subject)
		claims["family_name"] = fmt.Sprintf("%s-FamilyName", subject)
		claims["subject"] = user.Subject
	}

	code := this.generateRandomString()
	this.authParamsStore.addParams(code, authParams{
		IdTokenClaims: claims,
	})

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?state=%s&code=%s", c.Query("redirect_uri"),
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
