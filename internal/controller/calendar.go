package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/model"
	"github.com/shwatanap/gotion/internal/view/response"
)

func CalendarList(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID, _ := userIDAny.(string)
	o := model.NewGoogleOAuth()
	token, err := o.RefreshToken(c, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	cs, err := model.NewCalendarService(c, token)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	calendars, _ := cs.CalendarList()
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
