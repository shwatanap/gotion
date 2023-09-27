package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shwatanap/gotion/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// cookieからuser_idを取得
		user_id, err := c.Cookie("user_id")
		switch {
		case err == http.ErrNoCookie:
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
		}

		// userの存在確認
		err = model.IsUserExist(c, user_id)
		switch {
		case status.Code(err) == codes.NotFound:
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
		}

		c.Set("user_id", user_id)
	}
}
