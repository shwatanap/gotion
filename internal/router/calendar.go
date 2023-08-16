package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initCalendarRouter(router *gin.Engine) {
	router.GET("/calendars", controller.CalendarList)
}
