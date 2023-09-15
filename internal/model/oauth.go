package model

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shwatanap/gotion/internal/util"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OAuth struct {
	Config *oauth2.Config
}

func (o *OAuth) GetAuthCodeURL(oauthState string) string {
	// oauth2.ApprovalForce: ユーザーに強制的に認証を要求する
	authURL := o.Config.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return authURL
}

func (o *OAuth) RefreshToken(ctx context.Context, userID string) (*oauth2.Token, error) {
	refresh_token, err := GetRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	refresh_token, err = util.Decrypt([]byte(refresh_token), []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{
		RefreshToken: string(refresh_token),
	}
	// Token更新処理
	// TODO: Tokenが切れていた時の処理
	newToken, err := o.Config.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, err
	}
	cipherRefreshToken, err := util.Encrypt([]byte(newToken.RefreshToken), []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		return nil, err
	}
	if err := PutRefreshToken(ctx, userID, cipherRefreshToken); err != nil {
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

func GetRefreshToken(ctx context.Context, userID string) ([]byte, error) {
	client := NewFirestore(ctx)
	dsnap, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	m := dsnap.Data()
	return m["google_refresh_token"].([]byte), nil
}

func PutRefreshToken(ctx context.Context, userID string, refreshToken []byte) error {
	client := NewFirestore(ctx)
	docRef := client.Collection("users").Doc(userID)
	dsnap, err := docRef.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}
	// TODO: refresh_tokenの暗号化
	var data map[string]interface{}
	if dsnap.Exists() {
		data = dsnap.Data()
		data["google_refresh_token"] = refreshToken
		data["updated_at"] = time.Now()
	} else {
		data = map[string]interface{}{
			"user_id":              userID,
			"google_refresh_token": refreshToken,
			"created_at":           time.Now(),
			"updated_at":           time.Now(),
		}
	}
	_, err = docRef.Set(ctx, data)
	return err
}

// NotionはAccessTokenが永続的に使用できるので、RefreshTokenが存在しない
func PutNotionAccessToken(ctx context.Context, userID string, notionAccessToken []byte, dbURL string) error {
	client := NewFirestore(ctx)
	docRef := client.Collection("users").Doc(userID).Collection("connections").Doc(dbURL)
	// TODO: access_tokenの暗号化
	data := map[string]interface{}{
		"db_url":              dbURL,
		"notion_access_token": notionAccessToken,
		"created_at":          time.Now(),
		"updated_at":          time.Now(),
	}
	_, err := docRef.Set(ctx, data)
	return err
}
