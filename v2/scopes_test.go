package twitter

import (
	"reflect"
	"testing"
)

func TestScopeString(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []string
		expected string
	}{
		{
			name:     "single scope",
			scopes:   []string{ScopeTweetRead},
			expected: "tweet.read",
		},
		{
			name:     "multiple scopes",
			scopes:   []string{ScopeTweetRead, ScopeUsersRead, ScopeFollowsRead},
			expected: "tweet.read users.read follows.read",
		},
		{
			name:     "empty scopes",
			scopes:   []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScopeString(tt.scopes)
			if result != tt.expected {
				t.Errorf("ScopeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSplitScopes(t *testing.T) {
	tests := []struct {
		name        string
		scopeString string
		expected    []string
	}{
		{
			name:        "single scope",
			scopeString: "tweet.read",
			expected:    []string{"tweet.read"},
		},
		{
			name:        "multiple scopes",
			scopeString: "tweet.read users.read follows.read",
			expected:    []string{"tweet.read", "users.read", "follows.read"},
		},
		{
			name:        "empty string",
			scopeString: "",
			expected:    []string{},
		},
		{
			name:        "extra spaces",
			scopeString: "  tweet.read   users.read  ",
			expected:    []string{"tweet.read", "users.read"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitScopes(tt.scopeString)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitScopes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasScope(t *testing.T) {
	scopeString := "tweet.read users.read follows.read"
	
	tests := []struct {
		name        string
		targetScope string
		expected    bool
	}{
		{
			name:        "has scope",
			targetScope: "tweet.read",
			expected:    true,
		},
		{
			name:        "does not have scope",
			targetScope: "tweet.write",
			expected:    false,
		},
		{
			name:        "partial match should not work",
			targetScope: "tweet",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasScope(scopeString, tt.targetScope)
			if result != tt.expected {
				t.Errorf("HasScope() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateScopes(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []string
		expected []string
	}{
		{
			name:     "all valid scopes",
			scopes:   []string{ScopeTweetRead, ScopeUsersRead},
			expected: []string{},
		},
		{
			name:     "some invalid scopes",
			scopes:   []string{ScopeTweetRead, "invalid.scope", "another.invalid"},
			expected: []string{"invalid.scope", "another.invalid"},
		},
		{
			name:     "all invalid scopes",
			scopes:   []string{"invalid.scope", "another.invalid"},
			expected: []string{"invalid.scope", "another.invalid"},
		},
		{
			name:     "empty scopes",
			scopes:   []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateScopes(tt.scopes)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ValidateScopes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadOnlyScopes(t *testing.T) {
	// Test that ReadOnlyScopes contains only read scopes
	for _, scope := range ReadOnlyScopes {
		if !contains([]string{
			ScopeTweetRead, ScopeUsersRead, ScopeFollowsRead, ScopeLikeRead,
			ScopeListRead, ScopeBlockRead, ScopeMuteRead, ScopeSpaceRead, ScopeBookmarkRead,
		}, scope) {
			t.Errorf("ReadOnlyScopes contains non-read scope: %s", scope)
		}
	}
}

func TestWriteScopes(t *testing.T) {
	// Test that WriteScopes contains necessary read scopes for write operations
	expectedReads := []string{ScopeTweetRead, ScopeUsersRead, ScopeFollowsRead, ScopeLikeRead, ScopeListRead, ScopeBlockRead, ScopeMuteRead, ScopeBookmarkRead}
	
	for _, readScope := range expectedReads {
		if !contains(WriteScopes, readScope) {
			t.Errorf("WriteScopes missing necessary read scope: %s", readScope)
		}
	}
}

func TestScopeGroups(t *testing.T) {
	// Test that scope groups are not empty
	if len(ReadOnlyScopes) == 0 {
		t.Error("ReadOnlyScopes should not be empty")
	}
	
	if len(WriteScopes) == 0 {
		t.Error("WriteScopes should not be empty")
	}
	
	if len(AllScopes) == 0 {
		t.Error("AllScopes should not be empty")
	}
	
	if len(EssentialScopes) == 0 {
		t.Error("EssentialScopes should not be empty")
	}
	
	// Test that AllScopes contains more scopes than ReadOnlyScopes
	if len(AllScopes) <= len(ReadOnlyScopes) {
		t.Error("AllScopes should contain more scopes than ReadOnlyScopes")  
	}
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}