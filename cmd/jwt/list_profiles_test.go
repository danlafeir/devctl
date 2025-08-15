package jwt

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/danlafeir/devctl/pkg/secrets"
	"github.com/danlafeir/devctl/testutil/mocks"
)

func TestJWTListCommand_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "list", "--help")
	cmd.Dir = "../../"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "List all available OAuth client profiles") {
		t.Error("Help output does not contain expected description")
	}
}

func TestJWTListCommand_NoProfiles(t *testing.T) {
	mockSecrets := mocks.NewMockSecrets()
	secretsProvider = mockSecrets

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := jwtListCmd.RunE(jwtListCmd, nil)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	outputStr := buf.String()

	if err != nil {
		t.Fatalf("RunE returned error: %v", err)
	}
	if !strings.Contains(outputStr, "No profiles found.") {
		t.Errorf("Expected 'No profiles found.' message, got: %s", outputStr)
	}
}

func TestJWTListCommand_WithProfiles(t *testing.T) {
	mockSecrets := mocks.NewMockSecrets()
	secretsProvider = mockSecrets

	profiles := []string{"profile1", "profile2", "profile3"}
	for _, p := range profiles {
		mockSecrets.StoreOAuthClient(p, &secrets.OAuthClient{
			ClientID:     "id",
			ClientSecret: "secret",
			TokenURL:     "url",
			Scopes:       "scope",
			Audience:     "aud",
		})
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := jwtListCmd.RunE(jwtListCmd, nil)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	outputStr := buf.String()

	if err != nil {
		t.Fatalf("RunE returned error: %v", err)
	}
	for _, p := range profiles {
		if !strings.Contains(outputStr, p) {
			t.Errorf("Expected profile %q in output, got: %s", p, outputStr)
		}
	}
}
