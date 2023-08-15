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
	code := c.Query("code")
	// state検証
	if state != oauthState {
		log.Fatal("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/signup")
	}
	// Cookie削除
	c.SetCookie(OAUTH_STATE, oauthState, -1, "/oauth/google", "localhost", true, true)
	// Token保存
	o, _ := model.NewOAuth()
	token, _ := o.GetTokenFromCode(code)
	userId, err := o.GetUserId(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.SetCookie("user_id", userId, 365*24*60, "/", "localhost", true, true)
	if err = model.PutRefreshToken(c.Request.Context(), userId, token.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.Redirect(http.StatusFound, "/")
}
