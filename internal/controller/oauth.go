package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shwatanap/gotion/internal/model"
	"github.com/shwatanap/gotion/internal/util"
)

const GOOGLE_OAUTH_STATE = "google-oauth-state"
const NOTION_OAUTH_STATE = "notion-oauth-state"
const NOTION_ACCESS_TOKEN = "notion-access-token"

func GoogleSignUp(c *gin.Context) {
	id, _ := uuid.NewUUID()
	state := id.String()
	c.SetCookie(GOOGLE_OAUTH_STATE, state, 365*24*60, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
	o := model.NewGoogleOAuth()
	c.Header("Location", o.GetAuthCodeURL(state))
	c.JSON(http.StatusNoContent, gin.H{})
}

func GoogleSignUpCallback(c *gin.Context) {
	stateFromCookie, _ := c.Cookie(GOOGLE_OAUTH_STATE)
	stateFromRequest := c.Query("state")
	code := c.Query("code")
	// state検証
	if stateFromRequest != stateFromCookie {
		c.SetCookie(GOOGLE_OAUTH_STATE, "", -1, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid oauth google state",
		})
	}
	// Cookie削除
	c.SetCookie(GOOGLE_OAUTH_STATE, "", -1, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
	// Token保存
	o := model.NewGoogleOAuth()
	token, err := o.GetTokenFromCode(c.Request.Context(), code)
	if token.RefreshToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "refresh token is empty",
		})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	userID, err := c.Cookie("user_id")
	if err != nil {
		userID, err = o.GetUserID(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		c.SetCookie("user_id", userID, 365*24*60, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
	}
	cipherRefreshToken, err := util.Encrypt([]byte(token.RefreshToken), []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err = model.PutRefreshToken(c.Request.Context(), userID, cipherRefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.Redirect(http.StatusFound, "http://localhost:5173/step/notion-oauth")
}

func NotionOAuth(c *gin.Context) {
	id, _ := uuid.NewUUID()
	state := id.String()
	c.SetCookie(NOTION_OAUTH_STATE, state, 365*24*60, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
	o := model.NewNotionOAuth()
	c.Header("Location", o.GetAuthCodeURL(state))
	c.JSON(http.StatusNoContent, gin.H{})
}

func NotionOAuthCallback(c *gin.Context) {
	stateFromCookie, _ := c.Cookie(NOTION_OAUTH_STATE)
	stateFromRequest := c.Query("state")
	code := c.Query("code")
	// state検証
	if stateFromRequest != stateFromCookie {
		c.SetCookie(GOOGLE_OAUTH_STATE, "", -1, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid oauth google state",
		})
	}
	// Cookie削除
	c.SetCookie(NOTION_OAUTH_STATE, "", -1, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
	// Token保存
	o := model.NewNotionOAuth()
	token, err := o.GetTokenFromCode(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	cipherAccessToken, err := util.Encrypt([]byte(token.AccessToken), []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetCookie(NOTION_ACCESS_TOKEN, string(cipherAccessToken), 365*24*60, "/", os.Getenv("CLIENT_DOMAIN"), true, true)
	c.Redirect(http.StatusFound, "http://localhost:5173/step/input-db-name")
}
