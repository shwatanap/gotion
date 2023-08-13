package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("internal/view/*.html")
	initHealthRouter(router)
	initAuthRouter(router)

	router.GET("/", func(c *gin.Context) {
		// code := c.Query("code")
		// g, err := oauth.NewGoogle()
		// if err != nil {
		// 	// TODO: error handling
		// 	panic(err)
		// }

		// client, err := g.GetClient(code)
		// if err != nil {
		// 	// TODO: error handling
		// 	panic(err)
		// }

		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	return router
}
