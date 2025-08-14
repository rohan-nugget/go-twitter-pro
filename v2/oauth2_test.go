package twitter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewAuthorizerV2(t *testing.T) {
	auth := NewAuthorizerV2("access123", "refresh456", "client_id", "client_secret", nil)
	
	if auth.accessToken != "access123" {
		t.Errorf("Expected accessToken to be 'access123', got '%s'", auth.accessToken)
	}
	
	if auth.refreshToken != "refresh456" {
		t.Errorf("Expected refreshToken to be 'refresh456', got '%s'", auth.refreshToken)
	}
}

func TestAuthorizerV2_Add_BearerOnly(t *testing.T) {
	// Test bearer-only mode (no refresh token)
	auth := NewAuthorizerV2("test_token", "", "", "", nil)
	
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	auth.Add(req)
	
	authHeader := req.Header.Get("Authorization")
	expected := "Bearer test_token"
	
	if authHeader != expected {
		t.Errorf("Expected Authorization header to be '%s', got '%s'", expected, authHeader)
	}
}

func TestAuthorizerV2_RefreshToken(t *testing.T) {
	// Create a mock server for token refresh
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := OAuth2TokenResponse{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
			TokenType:    "Bearer",
			ExpiresIn:    7200,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	var callbackAccessToken, callbackRefreshToken string
	callback := func(accessToken, refreshToken string) {
		callbackAccessToken = accessToken
		callbackRefreshToken = refreshToken
	}
	
	auth := NewAuthorizerV2("old_access", "old_refresh", "client_id", "client_secret", callback)
	auth.tokenURL = server.URL
	
	// Force refresh
	ctx := context.Background()
	err := auth.ForceRefresh(ctx)
	if err != nil {
		t.Errorf("Expected force refresh to succeed, got error: %v", err)
	}
	
	// Verify tokens were updated
	access, refresh := auth.GetTokens()
	if access != "new_access_token" {
		t.Errorf("Expected access token to be 'new_access_token', got '%s'", access)
	}
	
	if refresh != "new_refresh_token" {
		t.Errorf("Expected refresh token to be 'new_refresh_token', got '%s'", refresh)
	}
	
	// Wait briefly for callback to execute (it runs in a goroutine)
	time.Sleep(100 * time.Millisecond)
	
	// Verify callback was called
	if callbackAccessToken != "new_access_token" {
		t.Errorf("Expected callback to receive 'new_access_token', got '%s'", callbackAccessToken)
	}
	_ = callbackRefreshToken // Avoid unused variable warning
}

func TestAuthorizerV2_NoRefreshToken(t *testing.T) {
	// Test that ForceRefresh fails when no refresh token
	auth := NewAuthorizerV2("access_token", "", "", "", nil)
	
	ctx := context.Background()
	err := auth.ForceRefresh(ctx)
	if err == nil {
		t.Error("Expected force refresh to fail without refresh token")
	}
	
	if !strings.Contains(err.Error(), "no refresh token available") {
		t.Errorf("Expected error about no refresh token, got: %v", err)
	}
}