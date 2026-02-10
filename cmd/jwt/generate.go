/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package jwt

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/danlafeir/devctl/pkg/config"
	"github.com/danlafeir/devctl/pkg/secrets"
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

	// Initialize config
	if err := config.InitConfig(""); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Get OAuth client config from config file
	profiles, err := config.FetchConfig("jwt")
	if err != nil {
		return fmt.Errorf("failed to fetch config: %w", err)
	}

	profileData, exists := profiles[profileFlag]
	if !exists {
		return fmt.Errorf("profile '%s' not found", profileFlag)
	}

	profileMap, ok := profileData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid profile configuration for '%s'", profileFlag)
	}

	// Extract config values
	clientID, _ := profileMap["client_id"].(string)
	clientSecretRef, _ := profileMap["client_secret"].(string)
	tokenURL, _ := profileMap["token_url"].(string)
	scopes, _ := profileMap["scopes"].(string)
	audience, _ := profileMap["audience"].(string)

	// Resolve client secret from secrets store if it's a reference
	var clientSecret string
	if strings.HasPrefix(clientSecretRef, "secret:") {
		secretToken := strings.TrimPrefix(clientSecretRef, "secret:")
		clientSecret, err = secrets.Read("jwt", secretToken)
		if err != nil {
			return fmt.Errorf("failed to read client secret from secrets store: %w", err)
		}
		clientSecret = strings.TrimSpace(clientSecret)
	} else {
		clientSecret = clientSecretRef
	}

	if clientID == "" || clientSecret == "" || tokenURL == "" {
		return fmt.Errorf("incomplete OAuth client configuration for profile '%s'", profileFlag)
	}

	// Create OAuth client
	client := &secrets.OAuthClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
		Audience:     audience,
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
