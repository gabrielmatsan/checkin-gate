package lib

import (
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

type GoogleOAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleOAuthProvider(clientID, clientSecret, redirectURL string) *GoogleOAuthProvider {
	return &GoogleOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

/**
 * Get the authorization URL for the Google OAuth flow
 * @param state - The state parameter to prevent CSRF attacks
 * @return The authorization URL
 */
func (p *GoogleOAuthProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *GoogleOAuthProvider) Enchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}


func (p *GoogleOAuthProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := p.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}



