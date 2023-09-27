package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, err := c.Cookie("user_id")
		switch {
		case err == http.ErrNoCookie:
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			c.Abort()
		}
		c.Set("user_id", user_id)
	}
}
