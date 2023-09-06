package model

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

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
