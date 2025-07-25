/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	profileFlag  string
	outputWriter io.Writer = os.Stdout
)

// jwtGenerateCmd represents the jwt generate command
var jwtGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a JWT token using configured OAuth client",
	Long: `Generate a JWT token using a configured OAuth client stored in the system secrets.
	
This command retrieves OAuth client credentials from the secrets using the specified profile.`,
	RunE: runJWTGenerate,
}

func init() {
	jwtCmd.AddCommand(jwtGenerateCmd)

	// Add the --profile flag
	jwtGenerateCmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "OAuth client profile name (required)")
	jwtGenerateCmd.MarkFlagRequired("profile")
}

func runJWTGenerate(cmd *cobra.Command, args []string) error {
	return runJWTGenerateWithWriter(cmd, args, outputWriter)
}

func runJWTGenerateWithWriter(cmd *cobra.Command, args []string, w io.Writer) error {
	if profileFlag == "" {
		return fmt.Errorf("profile flag is required")
	}

	// Get OAuth client from secrets (now via interface)
	client, err := secretsProvider.GetOAuthClient(profileFlag)
	if err != nil {
		return fmt.Errorf("failed to get OAuth client from secrets: %w", err)
	}

	// Use clientcredentials flow to get a token
	cfg := clientcredentials.Config{
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
		TokenURL:     client.TokenURL,
		Scopes:       nil,
	}
	if client.Scopes != "" {
		cfg.Scopes = append(cfg.Scopes, client.Scopes)
	}

	tok, err := cfg.Token(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get token from OAuth server: %w", err)
	}

	_, err = fmt.Fprintln(w, tok.AccessToken)
	return err
}
