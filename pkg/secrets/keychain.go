/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package secrets

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// OAuthClient represents an OAuth client configuration
type OAuthClient struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TokenURL     string `json:"token_url"`
	Scopes       string `json:"scopes"`
	Audience     string `json:"audience"`
}

// RealSecrets implements SecretsProvider using the system keychain
type RealSecrets struct{}

func (r *RealSecrets) GetOAuthClient(profile string) (*OAuthClient, error) {
	output, err := DefaultSecrets.Read(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve OAuth client from secrets: %w", err)
	}
	var client OAuthClient
	if err := json.Unmarshal([]byte(output), &client); err != nil {
		return nil, fmt.Errorf("failed to parse OAuth client data: %w", err)
	}
	return &client, nil
}

func (r *RealSecrets) StoreOAuthClient(profile string, client *OAuthClient) error {
	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("failed to serialize OAuth client: %w", err)
	}
	if err := DefaultSecrets.Write(profile, string(data)); err != nil {
		return fmt.Errorf("failed to store OAuth client in secrets: %w", err)
	}
	return nil
}

// TODO: Move ListOAuthProfiles to SecretsProvider abstraction in the future
func (r *RealSecrets) ListOAuthProfiles() ([]string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", "cli.devctl.oauth", "-g")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 44 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list OAuth profiles: %w", err)
	}
	lines := strings.Split(string(output), "\n")
	var profiles []string
	for _, line := range lines {
		if strings.Contains(line, `"acct"<blob>="`) {
			start := strings.Index(line, `"acct"<blob>="`) + 12
			end := strings.LastIndex(line, `"`)
			if start < end {
				profile := line[start:end]
				profiles = append(profiles, profile)
			}
		}
	}
	return profiles, nil
}

// TODO: Move DeleteOAuthClient to SecretsProvider abstraction in the future
func (r *RealSecrets) DeleteOAuthClient(profile string) error {
	cmd := exec.Command("security", "delete-generic-password", "-s", "cli.devctl.oauth", "-a", profile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete OAuth client from keychain: %w", err)
	}
	return nil
}

// SecretsAdapter abstracts secure storage for secrets (local or remote)
type SecretsAdapter interface {
	Write(key string, value string) error
	Read(key string) (string, error)
}

// MacOSSecretsAdapter implements SecretsAdapter using the MacOS keychain
type MacOSSecretsAdapter struct{}

func (m *MacOSSecretsAdapter) Write(key string, value string) error {
	// For now, treat key as the profile name and value as the JSON-encoded OAuthClient
	cmd := exec.Command("security", "add-generic-password", "-U", "-s", "cli.devctl.oauth", "-a", key, "-w", value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to store secret in keychain: %w", err)
	}
	return nil
}

func (m *MacOSSecretsAdapter) Read(key string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", "cli.devctl.oauth", "-a", key, "-w")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret from keychain: %w", err)
	}
	return string(output), nil
}

// DefaultSecrets is the global secrets adapter (defaults to MacOS keychain)
var DefaultSecrets SecretsAdapter = &MacOSSecretsAdapter{}
