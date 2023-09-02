package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shwatanap/gotion/internal/model"
)

func GCalendarExport(c *gin.Context) {
	userID, _ := c.Cookie("user_id")
	o := model.NewGoogleOAuth()
	token, err := o.RefreshToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	model.GCalendarExport(
		c.Request.Context(),
		token,
		"NotionAPIKey",
		"DBID",
		"Export Test",
	)
}
