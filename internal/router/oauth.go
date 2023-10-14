package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initGoogleOAuthRouter(r *gin.Engine) {
	grg := r.Group("/oauth/google")
	grg.GET("", controller.GoogleSignUp)
	grg.GET("/callback", controller.GoogleSignUpCallback)
}

func initNotionOAuthRouter(r *gin.RouterGroup) {
	nrg := r.Group("/oauth/notion")
	nrg.GET("", controller.NotionOAuth)
	nrg.GET("/callback", controller.NotionOAuthCallback)
}
