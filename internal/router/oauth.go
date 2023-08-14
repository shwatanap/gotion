package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initOAuthRouter(r *gin.Engine) {
	rg := r.Group("/oauth/google")
	rg.GET("/", controller.GoogleSignUp)
	rg.GET("/callback", controller.GoogleSignUpCallback)
}
