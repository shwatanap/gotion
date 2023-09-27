package model

import (
	"context"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
)

func NewVerifier(ctx context.Context) (*oidc.IDTokenVerifier, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}
	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	}
	verifier := provider.Verifier(oidcConfig)
	return verifier, nil
}
