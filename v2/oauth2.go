package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TokenRefreshCallback is called when tokens are refreshed
type TokenRefreshCallback func(accessToken, refreshToken string)

// OAuth2TokenResponse represents the response from Twitter's OAuth2 token endpoint
type OAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope,omitempty"`
}

// AuthorizerV2 implements the Authorizer interface with OAuth2 token refresh support
// If only access token is provided, works like simple bearer auth
// If refresh token is also provided, automatically refreshes tokens when needed
type AuthorizerV2 struct {
	mu           sync.RWMutex
	accessToken  string
	refreshToken string
	clientID     string
	clientSecret string
	tokenURL     string
	client       *http.Client
	callback     TokenRefreshCallback
	expiresAt    time.Time
}

// NewAuthorizerV2 creates a new authorizer with smart auth detection
// If refreshToken is empty: works as simple bearer auth (no refresh)
// If refreshToken is provided: automatically refreshes tokens when needed
func NewAuthorizerV2(accessToken, refreshToken, clientID, clientSecret string, callback TokenRefreshCallback) *AuthorizerV2 {
	auth := &AuthorizerV2{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     "https://api.twitter.com/2/oauth2/token",
		client:       http.DefaultClient,
		callback:     callback,
		expiresAt:    time.Now().Add(2 * time.Hour),
	}
	
	// If no refresh capability, don't need client credentials
	if refreshToken == "" {
		auth.clientID = ""
		auth.clientSecret = ""
		auth.tokenURL = ""
	}
	
	return auth
}

// Add implements the Authorizer interface and automatically refreshes tokens when needed
// If no refresh token is available, works like simple bearer auth
func (o *AuthorizerV2) Add(req *http.Request) {
	o.mu.RLock()
	
	// If no refresh token, just use access token (simple bearer auth mode)
	if o.refreshToken == "" {
		token := o.accessToken
		o.mu.RUnlock()
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return
	}
	o.mu.RUnlock()
	
	// OAuth2 mode with refresh capability
	token, err := o.getValidToken(req.Context())
	if err != nil {
		// If we can't refresh, use the current token and let the API return an error
		o.mu.RLock()
		token = o.accessToken
		o.mu.RUnlock()
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}

// getValidToken returns a valid access token, refreshing if necessary
func (o *AuthorizerV2) getValidToken(ctx context.Context) (string, error) {
	o.mu.RLock()
	
	// Check if token is still valid (with 5 minute buffer)
	if time.Now().Add(5*time.Minute).Before(o.expiresAt) {
		token := o.accessToken
		o.mu.RUnlock()
		return token, nil
	}
	
	refreshToken := o.refreshToken
	o.mu.RUnlock()
	
	// Token is expired or about to expire, refresh it
	return o.refreshAccessToken(ctx, refreshToken)
}

// refreshAccessToken refreshes the access token using the refresh token
func (o *AuthorizerV2) refreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", fmt.Errorf("no refresh token available")
	}
	
	// Prepare refresh request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", o.clientID)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create refresh request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(o.clientID, o.clientSecret)
	
	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		return "", fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, buf.String())
	}
	
	var tokenResp OAuth2TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}
	
	// Update tokens with lock
	o.mu.Lock()
	o.accessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		o.refreshToken = tokenResp.RefreshToken
	}
	
	// Calculate expiry time
	if tokenResp.ExpiresIn > 0 {
		o.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	} else {
		// Default to 2 hours if no expiry provided
		o.expiresAt = time.Now().Add(2 * time.Hour)
	}
	
	newAccessToken := o.accessToken
	newRefreshToken := o.refreshToken
	o.mu.Unlock()
	
	// Call callback if provided
	if o.callback != nil {
		go o.callback(newAccessToken, newRefreshToken)
	}
	
	return newAccessToken, nil
}

// UpdateTokens allows manual updating of tokens
func (o *AuthorizerV2) UpdateTokens(accessToken, refreshToken string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	if accessToken != "" {
		o.accessToken = accessToken
	}
	if refreshToken != "" {
		o.refreshToken = refreshToken
	}
}

// GetTokens returns the current tokens (useful for persistence)
func (o *AuthorizerV2) GetTokens() (accessToken, refreshToken string) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.accessToken, o.refreshToken
}

// ForceRefresh forces a token refresh regardless of expiry time
func (o *AuthorizerV2) ForceRefresh(ctx context.Context) error {
	o.mu.RLock()
	refreshToken := o.refreshToken
	o.mu.RUnlock()
	
	if refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}
	
	_, err := o.refreshAccessToken(ctx, refreshToken)
	return err
}