package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetSameSite(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.SetSameSite(http.SameSiteNoneMode)
		c.Next()
	})
}
