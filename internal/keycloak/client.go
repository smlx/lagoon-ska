package keycloak

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

// Client is a Keycloak admin client.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	jwtPubKey  *rsa.PublicKey
	log        *zap.Logger
}

// NewClient creates a new keycloak client.
func NewClient(ctx context.Context, log *zap.Logger, baseURL, username,
	password string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL %s: %v", baseURL, err)
	}
	httpClient, err := httpClient(ctx, *u, "master", username, password)
	if err != nil {
		return nil, fmt.Errorf("couldn't get keycloak http client: %v", err)
	}
	pubKey, err := publicKey(ctx, *u, "master")
	if err != nil {
		return nil, fmt.Errorf("couldn't get master realm public key: %v", err)
	}
	return &Client{
		baseURL:    u,
		httpClient: httpClient,
		jwtPubKey:  pubKey,
		log:        log,
	}, nil
}
