package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shwatanap/gotion/internal/model"
)

const OAUTH_STATE = "oauth-state"

func GoogleSignUp(c *gin.Context) {
	id, _ := uuid.NewUUID()
	oauthState := id.String()
	// TODO: 本番環境と開発環境でドメインを変える
	c.SetCookie(OAUTH_STATE, oauthState, 365*24*60, "/oauth/google", "localhost", true, true)
	o, _ := model.NewOAuth()
	c.Redirect(http.StatusFound, o.GetAuthCodeURL(oauthState))
}

func GoogleSignUpCallback(c *gin.Context) {
	oauthState, _ := c.Cookie(OAUTH_STATE)
	state := c.Query("state")
	// code := c.Query("code")

	// state検証
	if state != oauthState {
		log.Fatal("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/signup")
	}
	// Cookie削除
	c.SetCookie(OAUTH_STATE, oauthState, -1, "/oauth/google", "localhost", true, true)
	c.Redirect(http.StatusFound, "/")
}
