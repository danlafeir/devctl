package cmd

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/danlafeir/devctl/pkg/secrets"
)

func TestJWTDeleteCommand_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "delete", "--help")
	cmd.Dir = "../"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "Delete an OAuth client profile from secrets") {
		t.Error("Help output does not contain expected description")
	}
}

func TestJWTDeleteCommand_MissingArgument(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "delete")
	cmd.Dir = "../"
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error when running without profile argument, got nil")
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "accepts 1 arg(s), received 0") {
		t.Errorf("Error output does not mention missing argument: got: %s", outputStr)
	}
}

func TestJWTDeleteCommand_NonExistentProfile(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "delete", "non-existent-profile")
	cmd.Dir = "../"
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error when running with non-existent profile, got nil")
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "failed to delete profile") {
		t.Error("Error output does not mention deletion failure")
	}
}

func TestJWTDeleteCommand_Success_MockSecrets(t *testing.T) {
	mockSecrets := secrets.NewMockSecrets()
	secretsProvider = mockSecrets

	profile := "test-profile-delete"
	err := mockSecrets.StoreOAuthClient(profile, &secrets.OAuthClient{
		ClientID:     "id",
		ClientSecret: "secret",
		TokenURL:     "url",
		Scopes:       "scope",
		Audience:     "aud",
	})
	if err != nil {
		t.Fatalf("Failed to store mock OAuth client: %v", err)
	}

	err = mockSecrets.DeleteOAuthClient(profile)
	if err != nil {
		t.Fatalf("Failed to delete profile: %v", err)
	}

	_, err = mockSecrets.GetOAuthClient(profile)
	if err == nil {
		t.Error("Expected error when getting deleted profile, got nil")
	}
}
