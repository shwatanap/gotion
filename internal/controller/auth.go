package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/utils/oauth"
)

func RedirectOAuth(c *gin.Context) {
	g, err := oauth.NewGoogle()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "ok",
		})
	}
	c.Redirect(http.StatusFound, g.GetAuthCodeURL())
}
