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

// KeychainProvider defines the interface for keychain operations
// This allows for mocking in tests.
type KeychainProvider interface {
	GetOAuthClient(profile string) (*OAuthClient, error)
	StoreOAuthClient(profile string, client *OAuthClient) error
	ListOAuthProfiles() ([]string, error)
	DeleteOAuthClient(profile string) error
}

// realKeychain implements KeychainProvider using the system keychain
type realKeychain struct{}

func (r *realKeychain) GetOAuthClient(profile string) (*OAuthClient, error) {
	output, err := DefaultSecrets.Read(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve OAuth client from secrets: %w", err)
	}
	var client OAuthClient
	if err := json.Unmarshal(output, &client); err != nil {
		return nil, fmt.Errorf("failed to parse OAuth client data: %w", err)
	}
	return &client, nil
}

func (r *realKeychain) StoreOAuthClient(profile string, client *OAuthClient) error {
	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("failed to serialize OAuth client: %w", err)
	}
	if err := DefaultSecrets.Write(profile, data); err != nil {
		return fmt.Errorf("failed to store OAuth client in secrets: %w", err)
	}
	return nil
}

// TODO: Move ListOAuthProfiles to SecretsProvider abstraction in the future
func (r *realKeychain) ListOAuthProfiles() ([]string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", "devctl-oauth", "-g")
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
func (r *realKeychain) DeleteOAuthClient(profile string) error {
	cmd := exec.Command("security", "delete-generic-password", "-s", "devctl-oauth", "-a", profile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete OAuth client from keychain: %w", err)
	}
	return nil
}

// DefaultKeychain is the global keychain provider used in production
var DefaultKeychain KeychainProvider = &realKeychain{}

// For backward compatibility, keep the old function names as wrappers
func GetOAuthClient(profile string) (*OAuthClient, error) {
	return DefaultKeychain.GetOAuthClient(profile)
}
func StoreOAuthClient(profile string, client *OAuthClient) error {
	return DefaultKeychain.StoreOAuthClient(profile, client)
}
func ListOAuthProfiles() ([]string, error) {
	return DefaultKeychain.ListOAuthProfiles()
}
func DeleteOAuthClient(profile string) error {
	return DefaultKeychain.DeleteOAuthClient(profile)
}

// MockKeychain is an in-memory implementation for tests
// Not thread-safe, but sufficient for unit tests
type MockKeychain struct {
	store map[string]*OAuthClient
}

func NewMockKeychain() *MockKeychain {
	return &MockKeychain{store: make(map[string]*OAuthClient)}
}

func (m *MockKeychain) GetOAuthClient(profile string) (*OAuthClient, error) {
	c, ok := m.store[profile]
	if !ok {
		return nil, fmt.Errorf("profile not found")
	}
	return c, nil
}
func (m *MockKeychain) StoreOAuthClient(profile string, client *OAuthClient) error {
	m.store[profile] = client
	return nil
}
func (m *MockKeychain) ListOAuthProfiles() ([]string, error) {
	profiles := make([]string, 0, len(m.store))
	for k := range m.store {
		profiles = append(profiles, k)
	}
	return profiles, nil
}
func (m *MockKeychain) DeleteOAuthClient(profile string) error {
	delete(m.store, profile)
	return nil
}

// SecretsProvider abstracts secure storage for secrets (local or remote)
type SecretsProvider interface {
	Write(key string, value []byte) error
	Read(key string) ([]byte, error)
}

// MacOSKeychainAdapter implements SecretsProvider using the MacOS keychain
type MacOSKeychainAdapter struct{}

func (m *MacOSKeychainAdapter) Write(key string, value []byte) error {
	// For now, treat key as the profile name and value as the JSON-encoded OAuthClient
	cmd := exec.Command("security", "add-generic-password", "-U", "-s", "devctl-oauth", "-a", key, "-w", string(value))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to store secret in keychain: %w", err)
	}
	return nil
}

func (m *MacOSKeychainAdapter) Read(key string) ([]byte, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", "devctl-oauth", "-a", key, "-w")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret from keychain: %w", err)
	}
	return output, nil
}

// DefaultSecrets is the global secrets provider (defaults to MacOS keychain)
var DefaultSecrets SecretsProvider = &MacOSKeychainAdapter{}
