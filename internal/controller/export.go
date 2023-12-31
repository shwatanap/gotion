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
		return
	}
	userIDAny, _ := c.Get("user_id")
	userID, _ := userIDAny.(string)
	cipherAccessToken, _ := c.Cookie(NOTION_ACCESS_TOKEN)
	accessToken, err := util.Decrypt([]byte(cipherAccessToken), []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	o := model.NewGoogleOAuth()
	token, err := o.RefreshToken(c, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	db, err := model.GCalendarExport(
		c,
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
		return
	}
	if err := model.PutNotionAccessToken(c, userID, accessToken, db.ID.String(), req.DBName); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	res := response.ExportResponse{
		DBURL: db.URL,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}
