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
	o := model.NewGoogleOAuth()
	c.Redirect(http.StatusFound, o.GetAuthCodeURL(oauthState))
}

func GoogleSignUpCallback(c *gin.Context) {
	oauthState, _ := c.Cookie(OAUTH_STATE)
	state := c.Query("state")
	code := c.Query("code")
	// state検証
	if state != oauthState {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/signup")
	}
	// Cookie削除
	c.SetCookie(OAUTH_STATE, oauthState, -1, "/oauth/google", "localhost", true, true)
	// Token保存
	o := model.NewGoogleOAuth()
	token, err := o.GetTokenFromCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	userID, err := o.GetUserID(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.SetCookie("user_id", userID, 365*24*60, "/", "localhost", true, true)
	if err = model.PutRefreshToken(c.Request.Context(), userID, token.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.Redirect(http.StatusFound, "/")
}

func NotionOAuth(c *gin.Context) {
	id, _ := uuid.NewUUID()
	oauthState := id.String()
	// TODO: 本番環境と開発環境でドメインを変える
	c.SetCookie(OAUTH_STATE, oauthState, 365*24*60, "/oauth/notion", "localhost", true, true)
	o := model.NewNotionOAuth()
	c.Redirect(http.StatusFound, o.GetAuthCodeURL(oauthState))
}

func NotionOAuthCallback(c *gin.Context) {
	oauthState, _ := c.Cookie(OAUTH_STATE)
	state := c.Query("state")
	code := c.Query("code")
	// state検証
	if state != oauthState {
		log.Println("invalid oauth notion state")
		c.Redirect(http.StatusTemporaryRedirect, "/signup")
	}
	// Cookie削除
	c.SetCookie(OAUTH_STATE, oauthState, -1, "/oauth/notion", "localhost", true, true)
	// Token保存
	o := model.NewNotionOAuth()
	token, err := o.GetTokenFromCode(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("🥺", token, err)
	c.Redirect(http.StatusFound, "/")
}
