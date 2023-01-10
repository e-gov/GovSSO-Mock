package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *routeHandler) handleAuthCancel(c *gin.Context) {
	c.Redirect(http.StatusFound, fmt.Sprintf("%s?error=user_cancel&error_description=User+canceled+the+authentication+process.&state=%s",
		c.PostForm("redirect_uri"),
		c.PostForm("state")))
}
