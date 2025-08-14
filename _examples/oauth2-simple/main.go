package main

import (
	"fmt"
	"net/http"

	twitter "github.com/rohan-nugget/go-twitter-pro"
)

func main() {
	// Example 1: OAuth2 with auto-refresh
	fmt.Println("=== OAuth2 with auto-refresh ===")
	
	auth := twitter.NewAuthorizerV2(
		"access_token",
		"refresh_token", 
		"client_id",
		"client_secret",
		func(newAccess, newRefresh string) {
			fmt.Printf("Tokens refreshed!\n")
			saveTokens(newAccess, newRefresh)
		},
	)
	
	user := &twitter.User{
		Authorizer: auth,
		Client:     http.DefaultClient,
		Host:       "https://api.twitter.com",
	}
	
	fmt.Println("Created OAuth2 client with auto-refresh")
	
	// Example 2: Simple bearer token (no refresh)
	fmt.Println("\n=== Simple bearer token ===")
	
	bearerAuth := twitter.NewAuthorizerV2("access_token", "", "", "", nil)
	
	userSimple := &twitter.User{
		Authorizer: bearerAuth,
		Client:     http.DefaultClient,
		Host:       "https://api.twitter.com",
	}
	
	fmt.Println("Created simple bearer client")
	
	// Use the clients normally - tokens refresh automatically when needed
	_ = user
	_ = userSimple
}

func saveTokens(accessToken, refreshToken string) {
	// Implement your token storage logic here
	// Examples:
	// - Save to database
	// - Save to file  
	// - Update environment variables
	// - Send to external service
	fmt.Printf("Saving tokens to storage...\n")
}