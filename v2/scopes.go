package twitter

// Twitter API v2 OAuth2 Scopes
// Based on X (Twitter) API v2 documentation

// Individual Scopes
const (
	// Tweet Scopes
	ScopeTweetRead      = "tweet.read"       // Read Tweets
	ScopeTweetWrite     = "tweet.write"      // Create and delete Tweets
	ScopeTweetModerate  = "tweet.moderate"   // Hide and unhide replies to Tweets
	
	// User Scopes  
	ScopeUsersRead      = "users.read"       // Read user profile information
	ScopeUsersWrite     = "users.write"      // Update user profile
	
	// Follow Scopes
	ScopeFollowsRead    = "follows.read"     // Read follower/following lists
	ScopeFollowsWrite   = "follows.write"    // Follow and unfollow users
	
	// Like Scopes
	ScopeLikeRead       = "like.read"        // Read liked Tweets
	ScopeLikeWrite      = "like.write"       // Like and unlike Tweets
	
	// Retweet Scopes
	ScopeRetweetWrite   = "retweet.write"    // Retweet and undo retweets
	
	// List Scopes
	ScopeListRead       = "list.read"        // Read List information
	ScopeListWrite      = "list.write"       // Create and manage Lists
	
	// Block and Mute Scopes
	ScopeBlockRead      = "block.read"       // Read blocked users
	ScopeBlockWrite     = "block.write"      // Block and unblock users
	ScopeMuteRead       = "mute.read"        // Read muted users  
	ScopeMuteWrite      = "mute.write"       // Mute and unmute users
	
	// Space Scopes
	ScopeSpaceRead      = "space.read"       // Read Spaces information
	
	// Bookmark Scopes
	ScopeBookmarkRead   = "bookmark.read"    // Read bookmarked Tweets
	ScopeBookmarkWrite  = "bookmark.write"   // Bookmark and remove bookmarks
	
	// Direct Message Scopes
	ScopeDMRead         = "dm.read"          // Read Direct Messages
	ScopeDMWrite        = "dm.write"         // Send and manage Direct Messages
	
	// Offline Access
	ScopeOfflineAccess  = "offline.access"   // Maintain access when user not present
)

// Scope Groups for common use cases
var (
	// ReadOnlyScopes - All read-only permissions
	ReadOnlyScopes = []string{
		ScopeTweetRead,
		ScopeUsersRead,
		ScopeFollowsRead,
		ScopeLikeRead,
		ScopeListRead,
		ScopeBlockRead,
		ScopeMuteRead,
		ScopeSpaceRead,
		ScopeBookmarkRead,
		ScopeDMRead,
	}
	
	// WriteScopes - All write permissions (includes reads where necessary)
	WriteScopes = []string{
		ScopeTweetRead,
		ScopeTweetWrite,
		ScopeUsersRead,
		ScopeUsersWrite,
		ScopeFollowsRead,
		ScopeFollowsWrite,
		ScopeLikeRead,
		ScopeLikeWrite,
		ScopeRetweetWrite,
		ScopeListRead,
		ScopeListWrite,
		ScopeBlockRead,
		ScopeBlockWrite,
		ScopeMuteRead,
		ScopeMuteWrite,
		ScopeBookmarkRead,
		ScopeBookmarkWrite,
		ScopeDMRead,
		ScopeDMWrite,
	}
	
	// AllScopes - All available scopes including moderation and offline access
	AllScopes = []string{
		ScopeTweetRead,
		ScopeTweetWrite,
		ScopeTweetModerate,
		ScopeUsersRead,
		ScopeUsersWrite,
		ScopeFollowsRead,
		ScopeFollowsWrite,
		ScopeLikeRead,
		ScopeLikeWrite,
		ScopeRetweetWrite,
		ScopeListRead,
		ScopeListWrite,
		ScopeBlockRead,
		ScopeBlockWrite,
		ScopeMuteRead,
		ScopeMuteWrite,
		ScopeSpaceRead,
		ScopeBookmarkRead,
		ScopeBookmarkWrite,
		ScopeDMRead,
		ScopeDMWrite,
		ScopeOfflineAccess,
	}
	
	// Essential Scopes - Most commonly needed scopes
	EssentialScopes = []string{
		ScopeTweetRead,
		ScopeTweetWrite,
		ScopeUsersRead,
		ScopeFollowsRead,
		ScopeLikeRead,
		ScopeLikeWrite,
		ScopeOfflineAccess,
	}
	
	// Bot Scopes - Common scopes for bot applications
	BotScopes = []string{
		ScopeTweetRead,
		ScopeTweetWrite,
		ScopeUsersRead,
		ScopeFollowsRead,
		ScopeFollowsWrite,
		ScopeLikeRead,
		ScopeLikeWrite,
		ScopeRetweetWrite,
		ScopeOfflineAccess,
	}
	
	// Analytics Scopes - Read-only scopes for analytics applications
	AnalyticsScopes = []string{
		ScopeTweetRead,
		ScopeUsersRead,
		ScopeFollowsRead,
		ScopeLikeRead,
		ScopeListRead,
		ScopeSpaceRead,
	}
)

// ScopeString joins scopes with spaces (OAuth2 standard format)
func ScopeString(scopes []string) string {
	result := ""
	for i, scope := range scopes {
		if i > 0 {
			result += " "
		}
		result += scope
	}
	return result
}

// HasScope checks if a scope string contains a specific scope
func HasScope(scopeString, targetScope string) bool {
	scopes := SplitScopes(scopeString)
	for _, scope := range scopes {
		if scope == targetScope {
			return true
		}
	}
	return false
}

// SplitScopes splits a space-separated scope string into individual scopes
func SplitScopes(scopeString string) []string {
	if scopeString == "" {
		return []string{}
	}
	
	scopes := []string{}
	current := ""
	
	for _, char := range scopeString {
		if char == ' ' {
			if current != "" {
				scopes = append(scopes, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		scopes = append(scopes, current)
	}
	
	return scopes
}

// ValidateScopes checks if all provided scopes are valid Twitter API scopes
func ValidateScopes(scopes []string) []string {
	validScopes := make(map[string]bool)
	for _, scope := range AllScopes {
		validScopes[scope] = true
	}
	
	invalid := []string{}
	for _, scope := range scopes {
		if !validScopes[scope] {
			invalid = append(invalid, scope)
		}
	}
	
	return invalid
}