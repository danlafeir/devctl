/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package secrets

import (
	"fmt"

	"github.com/keybase/go-keychain"
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

// buildServiceName creates the service name using the naming convention: cli.devctl.<cmd>
func (r *RealSecrets) buildServiceName(cmd string) string {
	return fmt.Sprintf("cli.devctl.%s", cmd)
}

func (r *RealSecrets) Read(cmd, token string) (string, error) {
	serviceName := r.buildServiceName(cmd)
	
	data, err := keychain.GetGenericPassword(serviceName, token, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret from keychain: %w", err)
	}
	
	return string(data), nil
}

func (r *RealSecrets) Write(cmd, token, value string) error {
	serviceName := r.buildServiceName(cmd)
	
	item := keychain.NewGenericPassword(serviceName, token, "", []byte(value), "")
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	
	// Try to add the item
	err := keychain.AddItem(item)
	if err != nil {
		// If add fails (item might exist), try to update
		queryItem := keychain.NewGenericPassword(serviceName, token, "", nil, "")
		updateItem := keychain.NewGenericPassword(serviceName, token, "", []byte(value), "")
		err = keychain.UpdateItem(queryItem, updateItem)
		if err != nil {
			return fmt.Errorf("failed to store secret in keychain: %w", err)
		}
	}
	
	return nil
}

func (r *RealSecrets) List(cmd string) ([]string, error) {
	serviceName := r.buildServiceName(cmd)
	
	// Use the convenience function to get accounts for a service
	accounts, err := keychain.GetAccountsForService(serviceName)
	if err != nil {
		// If no items found, return empty list
		return []string{}, nil
	}
	
	return accounts, nil
}

func (r *RealSecrets) Delete(cmd, token string) error {
	serviceName := r.buildServiceName(cmd)
	
	err := keychain.DeleteGenericPasswordItem(serviceName, token)
	if err != nil {
		return fmt.Errorf("failed to delete secret from keychain: %w", err)
	}
	
	return nil
}

