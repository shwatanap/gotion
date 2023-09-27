package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shwatanap/gotion/internal/model"
	"github.com/shwatanap/gotion/internal/util"
)

// Cookie name
const (
	GOOGLE_OAUTH_STATE  = "google-oauth-state"
	GOOGLE_OAUTH_NONCE  = "google-oauth-nonce"
	NOTION_OAUTH_STATE  = "notion-oauth-state"
	NOTION_ACCESS_TOKEN = "notion-access-token"
)

// Cookie path
const (
	GOOGLE_OAUTH_PATH = "/oauth/google"
	NOTION_OAUTH_PATH = "/oauth/notion"
)

func GoogleSignUp(c *gin.Context) {
	state, _ := util.RandString(16)
	nonce, _ := util.RandString(16)
	util.SetCookie(c, GOOGLE_OAUTH_STATE, state, 365*24*60, GOOGLE_OAUTH_PATH, true, true)
	util.SetCookie(c, GOOGLE_OAUTH_NONCE, nonce, 365*24*60, GOOGLE_OAUTH_PATH, true, true)
	o := model.NewGoogleOAuth()
	c.Header("Location", o.GetAuthCodeURLWithNonce(state, nonce))
	c.JSON(http.StatusNoContent, gin.H{})
}

func GoogleSignUpCallback(c *gin.Context) {
	stateFromCookie, _ := c.Cookie(GOOGLE_OAUTH_STATE)
	stateFromRequest := c.Query("state")
	code := c.Query("code")
	// state検証
	if stateFromRequest != stateFromCookie {
		util.SetCookie(c, GOOGLE_OAUTH_STATE, "", -1, GOOGLE_OAUTH_PATH, true, true)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid oauth google state",
		})
		return
	}
	// Cookie削除
	util.SetCookie(c, GOOGLE_OAUTH_STATE, "", -1, GOOGLE_OAUTH_PATH, true, true)
	// Token保存
	o := model.NewGoogleOAuth()
	token, err := o.GetTokenFromCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if token == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "token is nil",
		})
		return
	}
	if token.RefreshToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "refresh token is empty",
		})
		return
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "id_token is empty",
		})
		return
	}
	verifier, _ := model.NewVerifier(c.Request.Context())
	// idTokenの検証と解析
	idToken, err := verifier.Verify(c, rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// nonce検証
	nonce, _ := c.Cookie(GOOGLE_OAUTH_NONCE)
	if idToken.Nonce != nonce {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid oauth google nonce",
		})
		return
	}
	// nonceのcookie削除
	util.SetCookie(c, GOOGLE_OAUTH_NONCE, "", -1, GOOGLE_OAUTH_PATH, true, true)
	userID, err := c.Cookie("user_id")
	if err != nil {
		util.SetCookie(c, "user_id", idToken.Subject, 365*24*60, GOOGLE_OAUTH_PATH, true, true)
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
		return
	}
	c.Redirect(http.StatusFound, os.Getenv("CLIENT_BASE_URL")+"/step/notion-oauth")
}

func NotionOAuth(c *gin.Context) {
	id, _ := uuid.NewUUID()
	state := id.String()
	util.SetCookie(c, NOTION_OAUTH_STATE, state, 365*24*60, NOTION_OAUTH_PATH, true, true)
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
		util.SetCookie(c, GOOGLE_OAUTH_STATE, "", -1, NOTION_OAUTH_PATH, true, true)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid oauth notion state",
		})
		return
	}
	// Cookie削除
	util.SetCookie(c, NOTION_OAUTH_STATE, "", -1, NOTION_OAUTH_PATH, true, true)
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
	util.SetCookie(c, NOTION_ACCESS_TOKEN, string(cipherAccessToken), 365*24*60, NOTION_OAUTH_PATH, true, true)
	c.Redirect(http.StatusFound, os.Getenv("CLIENT_BASE_URL")+"/step/input-db-name")
}
