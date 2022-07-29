package keycloak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"golang.org/x/oauth2"
)

func httpClient(ctx context.Context, u url.URL, realm, username,
	password string) (*http.Client, error) {
	u.Path = path.Join(u.Path,
		fmt.Sprintf("/auth/realms/%s/protocol/openid-connect/token", realm))
	config := &oauth2.Config{
		ClientID: "admin-cli",
		Endpoint: oauth2.Endpoint{
			TokenURL: u.String(),
			AuthURL:  u.String(),
		},
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Timeout: 10 * time.Second,
	})
	// authenticate for a token
	token, err := config.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("couldn't get token for credentials: %v", err)
	}
	// wrap the token in a *http.Client
	return config.Client(ctx, token), nil
}
