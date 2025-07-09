/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package keychain

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

// GetOAuthClient retrieves an OAuth client configuration from the keychain
func GetOAuthClient(profile string) (*OAuthClient, error) {
	// Use security command to get the OAuth client data from keychain
	cmd := exec.Command("security", "find-generic-password", "-s", "devctl-oauth", "-a", profile, "-w")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve OAuth client from keychain: %w", err)
	}

	// Parse the JSON data
	var client OAuthClient
	if err := json.Unmarshal(output, &client); err != nil {
		return nil, fmt.Errorf("failed to parse OAuth client data: %w", err)
	}

	return &client, nil
}

// StoreOAuthClient stores an OAuth client configuration in the keychain
func StoreOAuthClient(profile string, client *OAuthClient) error {
	// Serialize the client to JSON
	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("failed to serialize OAuth client: %w", err)
	}

	// Store in keychain using security command
	cmd := exec.Command("security", "add-generic-password", "-U", "-s", "devctl-oauth", "-a", profile, "-w", string(data))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to store OAuth client in keychain: %w", err)
	}

	return nil
}

// ListOAuthProfiles lists all available OAuth profiles in the keychain
func ListOAuthProfiles() ([]string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", "devctl-oauth", "-g")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 44 {
			// No profiles found
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list OAuth profiles: %w", err)
	}

	// Parse the output to extract account names (profiles)
	lines := strings.Split(string(output), "\n")
	var profiles []string
	
	for _, line := range lines {
		if strings.Contains(line, `"acct"<blob>="`) {
			// Extract the account name from the line
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

// DeleteOAuthClient removes an OAuth client configuration from the keychain
func DeleteOAuthClient(profile string) error {
	cmd := exec.Command("security", "delete-generic-password", "-s", "devctl-oauth", "-a", profile)
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete OAuth client from keychain: %w", err)
	}

	return nil
} 