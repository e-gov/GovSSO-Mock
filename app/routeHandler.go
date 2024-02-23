package main

import (
	"GovSSO-Mock/app/idtoken"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type routeHandler struct {
	config            config
	idTokenSigningKey *jose.SigningKey
	predefinedUsers   []user
	predefinedClients []client
	authParamsStore   *authParamsStore
	idTokenService    idtoken.Service
	httpClient        *http.Client
}

func (this *routeHandler) init() error {
	router := gin.Default()
	router.Static("/assets", "./ui/assets")
	router.LoadHTMLGlob("ui/template/*")
	router.Use(this.setDefaultHeaders)

	router.GET("/", this.displayHomePage)
	router.GET("/.well-known/openid-configuration", this.displayOpenIdConfiguration)
	router.GET("/.well-known/jwks.json", this.displayJwks)
	router.GET("/oauth2/auth", this.handleAuthGet)
	router.POST("/oauth2/auth", this.handleAuthPost)
	router.POST("/oauth2/cancel", this.handleAuthCancel)
	router.POST("/oauth2/token", this.handleAuthTokenGeneration)
	router.GET("/oauth2/sessions/logout", this.handleSessionLogout)
	router.POST("/backchannel/sessions/logout", this.handleBackchannelSessionLogout)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(this.authRequestValidation, authRequest{})
		v.RegisterAlias("valid_response_type", "required,eq=code")
		v.RegisterAlias("valid_scope", "required,oneof=openid 'openid phone' 'phone openid'")
		v.RegisterAlias("valid_state", "required,min=8")
	}

	return router.RunTLS(":"+this.config.ServerPort, this.config.TLSCertificate, this.config.TLSPrivateKey)
}

func (this *routeHandler) setDefaultHeaders(c *gin.Context) {
	c.Writer.Header().Set("X-Frame-Options", "deny")
	c.Writer.Header().Set("Content-Security-Policy", "connect-src 'self'; "+
		"default-src 'none'; "+
		"font-src 'self'; "+
		"img-src 'self' data:; "+
		"script-src 'unsafe-inline' 'self'; "+ // Inline scripts are allowed - unsafe
		"style-src 'self'; "+
		"base-uri 'none'; "+
		"frame-ancestors 'none'; "+
		"block-all-mixed-content")
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Writer.Header().Set("Expires", "0")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-XSS-Protection", "0")
}

func (this *routeHandler) displayHomePage(c *gin.Context) {
	logoutResultMessage := c.Query("logout-result-message")
	c.HTML(http.StatusOK, "home.html", gin.H{
		"PredefinedClients":   this.predefinedClients,
		"LogoutResultMessage": logoutResultMessage,
	})
}

func (this *routeHandler) displayOpenIdConfiguration(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.HTML(http.StatusOK, "openid-configuration.json", gin.H{
		"Host": this.config.HostUri(),
	})
}

func (this *routeHandler) displayJwks(c *gin.Context) {
	key := this.idTokenSigningKey.Key.(*rsa.PrivateKey)
	keys := &[]jose.JSONWebKey{{
		Algorithm: "RS256",
		Use:       "sig",
		Key:       &key.PublicKey,
		KeyID:     this.config.IdTokenSignKeyId,
	}}

	c.JSON(http.StatusOK, gin.H{
		"keys": keys,
	})
}

func (this *routeHandler) findFromPredefinedClients(clientId string) *client {
	for _, client := range this.predefinedClients {
		if client.ClientId == clientId {
			return &client
		}
	}
	return nil
}
