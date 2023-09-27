package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",             // 開発環境
			"https://gotion.vercel.app",         // 本番環境
			"https://gotion-staging.vercel.app", // staging環境
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"withCredentials",
		},
		ExposeHeaders: []string{
			"Location",
		},
		MaxAge:           24 * time.Hour,
		AllowCredentials: true,
	})
}
