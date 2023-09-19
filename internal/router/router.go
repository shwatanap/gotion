package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/middleware"
)

func Router() *gin.Engine {
	router := gin.Default()
	middleware.Cors(router)
	middleware.SetSameSite(router)

	initTemplateRouter(router)
	initHealthRouter(router)
	initOAuthRouter(router)
	initCalendarRouter(router)
	initExportRouter(router)
	return router
}
