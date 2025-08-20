/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package jwt

import (
	"fmt"
	"strings"

	"github.com/danlafeir/devctl/pkg/config"
	"github.com/danlafeir/devctl/pkg/secrets"
	"github.com/spf13/cobra"
)

// These are provided by main.go via cmd package
var BuildGitHash string

// jwtCmd represents the jwt command
var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "JWT token management",
	Long: `Manage JWT tokens including generation, validation, and OAuth client operations.
	
This command provides utilities for working with JWT tokens, including generating
tokens from configured OAuth clients stored in the system keychain.`,
}

// jwtDeleteCmd represents the jwt delete command (was delete-profile)
var jwtDeleteCmd = &cobra.Command{
	Use:   "delete [profile]",
	Short: "Delete an OAuth client profile from config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := args[0]
		
		// Initialize config
		if err := config.InitConfig(""); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		
		// Get existing profiles
		profiles, err := config.FetchConfig("jwt")
		if err != nil {
			return fmt.Errorf("failed to fetch config: %w", err)
		}
		
		// Check if profile exists
		profileData, exists := profiles[profile]
		if !exists {
			return fmt.Errorf("profile '%s' not found", profile)
		}
		
		// Delete associated secret if it exists
		if profileMap, ok := profileData.(map[string]interface{}); ok {
			if clientSecretRef, ok := profileMap["client_secret"].(string); ok {
				if strings.HasPrefix(clientSecretRef, "secret:") {
					secretToken := strings.TrimPrefix(clientSecretRef, "secret:")
					// Try to delete the secret using the secrets provider
					// We'll ignore errors since the secret might not exist
					_ = secrets.DefaultSecretsProvider.Delete("jwt", secretToken)
				}
			}
		}
		
		// Delete the profile using the config package
		if err := config.DeleteConfigValue("jwt", profile); err != nil {
			return fmt.Errorf("failed to delete profile from config: %w", err)
		}
		
		fmt.Printf("Profile '%s' deleted successfully.\n", profile)
		return nil
	},
}

// jwtListCmd represents the jwt list command (was list-profiles)
var jwtListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available OAuth client profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config
		if err := config.InitConfig(""); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		
		profiles, err := config.FetchConfig("jwt")
		if err != nil {
			return fmt.Errorf("failed to fetch config: %w", err)
		}
		
		if len(profiles) == 0 {
			fmt.Println("No profiles found.")
			return nil
		}
		
		for profileName := range profiles {
			fmt.Println(profileName)
		}
		return nil
	},
}

func GetJWTCommand() *cobra.Command {
	return jwtCmd
}

func init() {
	jwtCmd.AddCommand(jwtConfigureCmd)
	jwtCmd.AddCommand(jwtGenerateCmd)
	jwtCmd.AddCommand(jwtDeleteCmd)
	jwtCmd.AddCommand(jwtListCmd)
}
