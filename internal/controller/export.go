package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shwatanap/gotion/internal/model"
	"github.com/shwatanap/gotion/internal/util"
	"github.com/shwatanap/gotion/internal/view/request"
	"github.com/shwatanap/gotion/internal/view/response"
)

func GCalendarExport(c *gin.Context) {
	var req request.ExportRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	userID, _ := c.Cookie("user_id")
	cipherAccessToken, _ := c.Cookie(NOTION_ACCESS_TOKEN)
	accessToken, err := util.Decrypt([]byte(cipherAccessToken), []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	o := model.NewGoogleOAuth()
	token, err := o.RefreshToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	db, err := model.GCalendarExport(
		c.Request.Context(),
		token,
		string(accessToken),
		req.PageID,
		req.DBName,
		req.CalendarIDs,
	)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	if err := model.PutNotionAccessToken(c.Request.Context(), userID, accessToken, db.ID.String(), req.DBName); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	res := response.ExportResponse{
		DBURL: db.URL,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}
