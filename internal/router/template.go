package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initTemplateRouter(r *gin.Engine) {
	r.LoadHTMLGlob("internal/view/template/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})
}
