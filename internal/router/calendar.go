package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initCalendarRouter(router *gin.RouterGroup) {
	router.GET("/calendars", controller.CalendarList)
}
