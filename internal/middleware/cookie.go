package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetSameSite() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetSameSite(http.SameSiteNoneMode)
		c.Next()
	}
}
