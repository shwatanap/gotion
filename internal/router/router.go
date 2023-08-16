package router

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	initTemplateRouter(router)
	initHealthRouter(router)
	initOAuthRouter(router)
	initCalendarRouter(router)
	return router
}
