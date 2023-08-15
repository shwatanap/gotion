package model

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type OAuth struct {
	Config *oauth2.Config
}

func NewOAuth() (*OAuth, error) {
	b, err := os.ReadFile("credentials_web.json")
	if err != nil {
		return nil, fmt.Errorf("nable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("nable to parse client secret file to config: %v", err)
	}
	oauth := &OAuth{
		Config: config,
	}
	return oauth, nil
}

func (g *OAuth) GetAuthCodeURL(oauthState string) string {
	// TODO: stateをどうするか
	authURL := g.Config.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
	return authURL
}

func (g *OAuth) GetClient(authCode string) (*http.Client, error) {
	token, err := g.Config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	// TODO: c.Request.Context()を使うべきか？
	return g.Config.Client(context.Background(), token), nil
}
