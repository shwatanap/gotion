package model

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OAuth struct {
	Config *oauth2.Config
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
