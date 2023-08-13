package oauth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Google struct {
	Config *oauth2.Config
}

func NewGoogle() (*Google, error) {
	return newGoogle()
}

func newGoogle() (*Google, error) {
	b, err := os.ReadFile("credentials_web.json")
	if err != nil {
		return nil, fmt.Errorf("nable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("nable to parse client secret file to config: %v", err)
	}
	google := &Google{
		Config: config,
	}
	return google, nil
}

func (g *Google) GetAuthCodeURL() string {
	// TODO: stateをどうするか
	authURL := g.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

func (g *Google) GetClient(authCode string) (*http.Client, error) {
	token, err := g.Config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	// TODO: c.Request.Context()を使うべきか？
	return g.Config.Client(context.Background(), token), nil
}
