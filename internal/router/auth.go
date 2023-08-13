package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initAuthRouter(r *gin.Engine) {
	r.GET("/auth", controller.RedirectOAuth)
}
