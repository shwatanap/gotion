package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shwatanap/gotion/internal/controller"
)

func initOAuthRouter(r *gin.Engine) {
	grg := r.Group("/oauth/google")
	grg.GET("/", controller.GoogleSignUp)
	grg.GET("/callback", controller.GoogleSignUpCallback)

	nrg := r.Group("/oauth/notion")
	nrg.GET("/", controller.NotionOAuth)
	nrg.GET("/callback", controller.NotionOAuthCallback)
}
