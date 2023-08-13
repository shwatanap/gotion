package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initHealthRouter(router *gin.Engine) {
	router.GET("/health", controller.Health)
}
