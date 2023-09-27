package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/middleware"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Cors())
	router.Use(middleware.SetSameSite())

	initTemplateRouter(router)
	initHealthRouter(router)

	authRequiredGroup := router.Group("/")
	authRequiredGroup.Use(middleware.Auth())
	{
		initOAuthRouter(authRequiredGroup)
		initCalendarRouter(authRequiredGroup)
		initExportRouter(authRequiredGroup)
	}
	return router
}
