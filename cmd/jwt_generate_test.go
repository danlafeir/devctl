/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/danlafeir/devctl/pkg/keychain"
)

func TestJWTGenerateCommand_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "generate", "--help")
	cmd.Dir = "../"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Generate a JWT token using a configured OAuth client stored in the system keychain.") {
		t.Error("Help output does not contain expected description")
	}
	if !strings.Contains(outputStr, "--profile") {
		t.Error("Help output does not contain --profile flag")
	}
}

func TestJWTGenerateCommand_MissingProfile(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "generate")
	cmd.Dir = "../"
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error when running without profile flag, got nil")
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "required flag(s) \"profile\" not set") {
		t.Errorf("Error output does not mention required flag: got: %s", outputStr)
	}
}

func TestJWTGenerateCommand_NonExistentProfile(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "generate", "--profile", "non-existent-profile")
	cmd.Dir = "../"
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error when running with non-existent profile, got nil")
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "failed to get OAuth client from keychain") {
		t.Error("Error output does not mention keychain retrieval failure")
	}
}

// The following tests require a real OAuth2 server to return a valid access token.
// If you want to run them, set up a test OAuth2 server and update the test profile accordingly.
// For now, we check that the command runs and returns a non-empty string (token or error).

// func TestJWTGenerateCommand_ValidProfile(t *testing.T) { ... }
// func TestJWTGenerateCommand_ShortFlag(t *testing.T) { ... }
// func TestJWTGenerateCommand_InvalidPrivateKey(t *testing.T) { ... }

// isMacOS checks if the current platform is macOS
func isMacOS() bool {
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "Darwin"
}

// mockOAuthServer starts a simple OAuth2 token endpoint for testing
func mockOAuthServer(t *testing.T, tokenValue string) (serverURL string, closeFn func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"access_token": "` + tokenValue + `", "token_type": "bearer", "expires_in": 3600}`))
			return
		}
		w.WriteHeader(404)
	}))
	return ts.URL + "/token", ts.Close
}

func TestJWTGenerateCommand_ValidProfile_MockOAuth(t *testing.T) {
	// Use a mock keychain for this test
	mockKC := keychain.NewMockKeychain()
	keychainProvider = mockKC

	tokenValue := "mocked-token-123"
	tokenURL, closeServer := mockOAuthServer(t, tokenValue)
	defer closeServer()

	profile := "test-profile-mock-oauth"
	// Store a mock OAuth client in the mock keychain
	err := mockKC.StoreOAuthClient(profile, &keychain.OAuthClient{
		ClientID:     "id",
		ClientSecret: "secret",
		TokenURL:     tokenURL,
		Scopes:       "scope1",
		Audience:     "aud",
	})
	if err != nil {
		t.Fatalf("Failed to store mock OAuth client: %v", err)
	}

	// Simulate the CLI logic directly by calling runJWTGenerateWithWriter
	profileFlag = profile
	var sb strings.Builder
	err = runJWTGenerateWithWriter(nil, nil, &sb)
	if err != nil {
		t.Fatalf("runJWTGenerateWithWriter failed: %v", err)
	}
	outStr := strings.TrimSpace(sb.String())
	if outStr != tokenValue {
		t.Errorf("Expected token %q, got %q", tokenValue, outStr)
	}
}
