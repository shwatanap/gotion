package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/model"
	"github.com/shwatanap/gotion/internal/view/response"
)

func CalendarList(c *gin.Context) {
	// TODO: Cookieにuser_idが存在しない場合のエラーハンドリング
	userID, _ := c.Cookie("user_id")
	o, _ := model.NewOAuth()
	token, err := o.RefreshToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	cs, err := model.NewCalendarService(c.Request.Context(), token)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	calendars, _ := cs.List()
	// var res model.CalendarListResponse
	res := make([]response.CalendarResponse, len(calendars))
	for i, c := range calendars {
		res[i] = response.CalendarResponse{
			ID:         c.Calendar.Id,
			Summary:    c.Calendar.Summary,
			ColorID:    c.Calendar.ColorId,
			AccessRole: c.Calendar.AccessRole,
		}
	}
	c.JSON(200, gin.H{
		"message": response.CalendarListResponse{Calendars: res},
	})
}
