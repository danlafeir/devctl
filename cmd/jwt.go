/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// jwtCmd represents the jwt command
var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "JWT token management",
	Long: `Manage JWT tokens including generation, validation, and OAuth client operations.
	
This command provides utilities for working with JWT tokens, including generating
tokens from configured OAuth clients stored in the system keychain.`,
}

func init() {
	rootCmd.AddCommand(jwtCmd)
} 