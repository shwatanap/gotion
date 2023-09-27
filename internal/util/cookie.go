package util

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name, value string, maxAge int, path string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		SameSite: http.SameSiteNoneMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}
