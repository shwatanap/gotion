package model

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OAuth struct {
	Config *oauth2.Config
}

func NewGoogleOAuth() *OAuth {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("GOOGLE_AUTH_URL"),
			TokenURL: os.Getenv("GOOGLE_TOKEN_URL"),
		},
		RedirectURL: os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			goauth2.UserinfoProfileScope,
			calendar.CalendarReadonlyScope,
		},
	}
	oauth := &OAuth{
		Config: cfg,
	}
	return oauth
}

func NewNotionOAuth() *OAuth {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("NOTION_CLIENT_ID"),
		ClientSecret: os.Getenv("NOTION_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("NOTION_AUTH_URL"),
			TokenURL: os.Getenv("NOTION_TOKEN_URL"),
		},
		RedirectURL: os.Getenv("NOTION_REDIRECT_URL"),
	}
	oauth := &OAuth{
		Config: cfg,
	}
	return oauth
}

func (o *OAuth) GetAuthCodeURL(oauthState string) string {
	authURL := o.Config.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
	return authURL
}

func (o *OAuth) RefreshToken(ctx context.Context, userID string) (*oauth2.Token, error) {
	refresh_token, _ := GetRefreshToken(ctx, userID)
	token := &oauth2.Token{
		RefreshToken: refresh_token,
	}
	// Token更新処理
	// TODO: Tokenが切れていた時の処理
	newToken, err := o.Config.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, err
	}
	if err := PutRefreshToken(ctx, userID, newToken.RefreshToken); err != nil {
		return nil, err
	}
	return newToken, nil
}

// 別のModelに書きべきかも
// Why: oauth2とは別のUserinfoAPIを叩くため
func (o *OAuth) GetUserID(ctx context.Context, token *oauth2.Token) (string, error) {
	os, err := goauth2.NewService(ctx, option.WithTokenSource(o.Config.TokenSource(ctx, token)))
	if err != nil {
		return "", err
	}
	us := goauth2.NewUserinfoService(os)
	userinfo, err := us.Get().Do()
	if err != nil {
		return "", err
	}
	return userinfo.Id, nil
}

func (o *OAuth) GetTokenFromCode(ctx context.Context, authCode string) (*oauth2.Token, error) {
	token, err := o.Config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("excahnge code: %w", err)
	}
	return token, nil
}

func GetRefreshToken(ctx context.Context, userID string) (string, error) {
	client := NewFirestore(ctx)
	dsnap, err := client.Collection("tokens").Doc(userID).Get(ctx)
	if err != nil {
		return "", err
	}
	m := dsnap.Data()
	return m["refresh_token"].(string), nil
}

func PutRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	client := NewFirestore(ctx)
	docRef := client.Collection("tokens").Doc(userID)
	dsnap, err := docRef.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}
	// TODO: refresh_tokenの暗号化
	data := map[string]interface{}{
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	}
	if !dsnap.Exists() {
		data["user_id"] = userID
		data["created_at"] = time.Now()
	}
	_, err = docRef.Set(ctx, data)
	return err
}
