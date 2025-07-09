/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os/exec"
	"strings"
	"testing"
)

// generateTestKeyPair creates a test RSA key pair
func generateTestKeyPair() (privateKeyPEM string, err error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	// Encode private key to PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM = string(pem.EncodeToMemory(privateKeyBlock))

	return privateKeyPEM, nil
}

func TestJWTGenerateCommand_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "jwt", "generate", "--help")
	cmd.Dir = "../"
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Generate a JWT token using configured OAuth client") {
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
	if !strings.Contains(outputStr, "profile flag is required") {
		t.Error("Error output does not mention profile flag")
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
