package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shwatanap/gotion/internal/model"
)

const GOOGLE_OAUTH_STATE = "google-oauth-state"
const NOTION_OAUTH_STATE = "notion-oauth-state"

func GoogleSignUp(c *gin.Context) {
	id, _ := uuid.NewUUID()
	oauthState := id.String()
	// TODO: 本番環境と開発環境でドメインを変える
	c.SetCookie(GOOGLE_OAUTH_STATE, oauthState, 365*24*60, "/", "localhost", true, true)
	o := model.NewGoogleOAuth()
	c.Header("Location", o.GetAuthCodeURL(oauthState))
	c.JSON(http.StatusFound, gin.H{})
}

func GoogleSignUpCallback(c *gin.Context) {
	oauthState, _ := c.Cookie(GOOGLE_OAUTH_STATE)
	state := c.Query("state")
	code := c.Query("code")
	// state検証
	if state != oauthState {
		c.SetCookie(GOOGLE_OAUTH_STATE, oauthState, -1, "/", "localhost", true, true)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid oauth google state",
		})
	}
	// Cookie削除
	c.SetCookie(GOOGLE_OAUTH_STATE, oauthState, -1, "/", "localhost", true, true)
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
	c.Redirect(http.StatusFound, "http://localhost:5173/next")
}

func NotionOAuth(c *gin.Context) {
	id, _ := uuid.NewUUID()
	oauthState := id.String()
	// TODO: 本番環境と開発環境でドメインを変える
	c.SetCookie(NOTION_OAUTH_STATE, oauthState, 365*24*60, "/", "localhost", true, true)
	o := model.NewNotionOAuth()
	c.Header("Location", o.GetAuthCodeURL(oauthState))
	c.JSON(http.StatusFound, gin.H{})
}

func NotionOAuthCallback(c *gin.Context) {
	oauthState, _ := c.Cookie(NOTION_OAUTH_STATE)
	state := c.Query("state")
	code := c.Query("code")
	// state検証
	if state != oauthState {
		c.SetCookie(GOOGLE_OAUTH_STATE, oauthState, -1, "/", "localhost", true, true)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid oauth google state",
		})
	}
	// Cookie削除
	c.SetCookie(NOTION_OAUTH_STATE, oauthState, -1, "/", "localhost", true, true)
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
	c.Redirect(http.StatusFound, "http://localhost:5173/next")
}
