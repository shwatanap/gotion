package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initExportRouter(router *gin.Engine) {
	router.POST("/export", controller.GCalendarExport)
}
