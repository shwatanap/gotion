package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shwatanap/gotion/internal/model"
)

func GCalendarExport(c *gin.Context) {
	userId, _ := c.Cookie("user_id")
	o, _ := model.NewOAuth()
	token, err := o.RefreshToken(c.Request.Context(), userId)
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
