/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package keychain

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
)

func TestOAuthClient_JSON(t *testing.T) {
	client := &OAuthClient{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     "https://example.com/oauth/token",
		Scopes:       "read write",
		Audience:     "https://api.example.com",
	}

	// Test marshaling
	data, err := json.Marshal(client)
	if err != nil {
		t.Fatalf("Failed to marshal OAuth client: %v", err)
	}

	// Test unmarshaling
	var unmarshaledClient OAuthClient
	err = json.Unmarshal(data, &unmarshaledClient)
	if err != nil {
		t.Fatalf("Failed to unmarshal OAuth client: %v", err)
	}

	// Verify fields match
	if unmarshaledClient.ClientID != client.ClientID {
		t.Errorf("ClientID mismatch: got %s, want %s", unmarshaledClient.ClientID, client.ClientID)
	}
	if unmarshaledClient.ClientSecret != client.ClientSecret {
		t.Errorf("ClientSecret mismatch: got %s, want %s", unmarshaledClient.ClientSecret, client.ClientSecret)
	}
	if unmarshaledClient.TokenURL != client.TokenURL {
		t.Errorf("TokenURL mismatch: got %s, want %s", unmarshaledClient.TokenURL, client.TokenURL)
	}
	if unmarshaledClient.Scopes != client.Scopes {
		t.Errorf("Scopes mismatch: got %s, want %s", unmarshaledClient.Scopes, client.Scopes)
	}
	if unmarshaledClient.Audience != client.Audience {
		t.Errorf("Audience mismatch: got %s, want %s", unmarshaledClient.Audience, client.Audience)
	}
}

func TestGetOAuthClient_NonExistent(t *testing.T) {
	// Test getting a non-existent profile
	_, err := GetOAuthClient("non-existent-profile")
	if err == nil {
		t.Error("Expected error when getting non-existent profile, got nil")
	}
}

func TestStoreAndGetOAuthClient(t *testing.T) {
	// Skip if not on macOS (keychain only works on macOS)
	if !isMacOS() {
		t.Skip("Skipping keychain tests on non-macOS platform")
	}

	testProfile := "test-profile-jwt-generate"
	testClient := &OAuthClient{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     "https://example.com/oauth/token",
		Scopes:       "read write",
		Audience:     "https://api.example.com",
	}

	// Clean up any existing test data
	DeleteOAuthClient(testProfile)

	// Test storing OAuth client
	err := StoreOAuthClient(testProfile, testClient)
	if err != nil {
		t.Fatalf("Failed to store OAuth client: %v", err)
	}

	// Test retrieving OAuth client
	retrievedClient, err := GetOAuthClient(testProfile)
	if err != nil {
		t.Fatalf("Failed to get OAuth client: %v", err)
	}

	// Verify retrieved client matches original
	if retrievedClient.ClientID != testClient.ClientID {
		t.Errorf("ClientID mismatch: got %s, want %s", retrievedClient.ClientID, testClient.ClientID)
	}
	if retrievedClient.ClientSecret != testClient.ClientSecret {
		t.Errorf("ClientSecret mismatch: got %s, want %s", retrievedClient.ClientSecret, testClient.ClientSecret)
	}
	if retrievedClient.TokenURL != testClient.TokenURL {
		t.Errorf("TokenURL mismatch: got %s, want %s", retrievedClient.TokenURL, testClient.TokenURL)
	}
	if retrievedClient.Scopes != testClient.Scopes {
		t.Errorf("Scopes mismatch: got %s, want %s", retrievedClient.Scopes, testClient.Scopes)
	}
	if retrievedClient.Audience != testClient.Audience {
		t.Errorf("Audience mismatch: got %s, want %s", retrievedClient.Audience, testClient.Audience)
	}

	// Clean up
	DeleteOAuthClient(testProfile)
}

func TestDeleteOAuthClient(t *testing.T) {
	// Skip if not on macOS
	if !isMacOS() {
		t.Skip("Skipping keychain tests on non-macOS platform")
	}

	testProfile := "test-profile-delete"
	testClient := &OAuthClient{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     "https://example.com/oauth/token",
		Scopes:       "read write",
		Audience:     "https://api.example.com",
	}

	// Store a test client
	err := StoreOAuthClient(testProfile, testClient)
	if err != nil {
		t.Fatalf("Failed to store OAuth client: %v", err)
	}

	// Verify it exists
	_, err = GetOAuthClient(testProfile)
	if err != nil {
		t.Fatalf("Failed to get OAuth client before deletion: %v", err)
	}

	// Delete the client
	err = DeleteOAuthClient(testProfile)
	if err != nil {
		t.Fatalf("Failed to delete OAuth client: %v", err)
	}

	// Verify it's gone
	_, err = GetOAuthClient(testProfile)
	if err == nil {
		t.Error("Expected error when getting deleted profile, got nil")
	}
}

func TestListOAuthProfiles(t *testing.T) {
	// Skip if not on macOS
	if !isMacOS() {
		t.Skip("Skipping keychain tests on non-macOS platform")
	}

	profiles, err := ListOAuthProfiles()
	if err != nil {
		t.Fatalf("Failed to list OAuth profiles: %v", err)
	}

	// Should return a slice (might be empty if no profiles exist)
	if profiles == nil {
		t.Error("Expected profiles to be a slice, got nil")
	}
}

// isMacOS checks if the current platform is macOS
func isMacOS() bool {
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "Darwin"
} 