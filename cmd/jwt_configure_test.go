package cmd

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/danlafeir/devctl/pkg/keychain"
)

func TestJWTConfigureCommand_Flags(t *testing.T) {
	testProfile := "test-profile-configure-flags"

	// Clean up before and after
	keychain.DeleteOAuthClient(testProfile)
	defer keychain.DeleteOAuthClient(testProfile)

	cmd := exec.Command("go", "run", "main.go", "jwt", "configure",
		"--profile", testProfile,
		"--client-id", "id-flags",
		"--client-secret", "secret-flags",
		"--token-url", "https://token-flags",
		"--scopes", "scope1 scope2",
		"--audience", "aud-flags",
	)
	cmd.Dir = "../"
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run jwt configure with flags: %v\nOutput: %s", err, string(output))
	}

	// Check that the profile was stored
	client, err := keychain.GetOAuthClient(testProfile)
	if err != nil {
		t.Fatalf("Failed to get stored profile: %v", err)
	}
	if client.ClientID != "id-flags" {
		t.Errorf("ClientID mismatch: got %s, want id-flags", client.ClientID)
	}
	if client.ClientSecret != "secret-flags" {
		t.Errorf("ClientSecret mismatch: got %s, want secret-flags", client.ClientSecret)
	}
	if client.TokenURL != "https://token-flags" {
		t.Errorf("TokenURL mismatch: got %s, want https://token-flags", client.TokenURL)
	}
	if client.Scopes != "scope1 scope2" {
		t.Errorf("Scopes mismatch: got %s, want scope1 scope2", client.Scopes)
	}
	if client.Audience != "aud-flags" {
		t.Errorf("Audience mismatch: got %s, want aud-flags", client.Audience)
	}
}

func TestJWTConfigureCommand_Interactive(t *testing.T) {
	testProfile := "test-profile-configure-interactive"

	// Clean up before and after
	keychain.DeleteOAuthClient(testProfile)
	defer keychain.DeleteOAuthClient(testProfile)

	// Simulate user input for all prompts
	input := strings.Join([]string{
		testProfile,
		"id-interactive",
		"secret-interactive",
		"https://token-interactive",
		"scopeA scopeB",
		"aud-interactive",
	}, "\n") + "\n"

	cmd := exec.Command("go", "run", "main.go", "jwt", "configure")
	cmd.Stdin = strings.NewReader(input)
	cmd.Dir = "../"
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run jwt configure interactively: %v\nOutput: %s", err, string(output))
	}

	// Check that the profile was stored
	client, err := keychain.GetOAuthClient(testProfile)
	if err != nil {
		t.Fatalf("Failed to get stored profile: %v", err)
	}
	if client.ClientID != "id-interactive" {
		t.Errorf("ClientID mismatch: got %s, want id-interactive", client.ClientID)
	}
	if client.ClientSecret != "secret-interactive" {
		t.Errorf("ClientSecret mismatch: got %s, want secret-interactive", client.ClientSecret)
	}
	if client.TokenURL != "https://token-interactive" {
		t.Errorf("TokenURL mismatch: got %s, want https://token-interactive", client.TokenURL)
	}
	if client.Scopes != "scopeA scopeB" {
		t.Errorf("Scopes mismatch: got %s, want scopeA scopeB", client.Scopes)
	}
	if client.Audience != "aud-interactive" {
		t.Errorf("Audience mismatch: got %s, want aud-interactive", client.Audience)
	}
} 