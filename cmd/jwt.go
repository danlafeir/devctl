/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/danlafeir/devctl/pkg/secrets"
	"github.com/spf13/cobra"
)

// Add this to ensure secretsProvider is shared across files
var secretsProvider secrets.SecretsProvider = secrets.DefaultSecretsProvider

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
	Short: "Delete an OAuth client profile from secrets",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := args[0]
		err := secretsProvider.DeleteOAuthClient(profile)
		if err != nil {
			return fmt.Errorf("failed to delete profile '%s': %w", profile, err)
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
		profiles, err := secretsProvider.ListOAuthProfiles()
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}
		if len(profiles) == 0 {
			fmt.Println("No profiles found.")
			return nil
		}
		for _, profile := range profiles {
			fmt.Println(profile)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(jwtCmd)
	jwtCmd.AddCommand(jwtDeleteCmd)
	jwtCmd.AddCommand(jwtListCmd)
}
