package jwt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/danlafeir/devctl/pkg/secrets"
	"github.com/spf13/cobra"
)

var (
	cfgProfile               string
	cfgClientID              string
	cfgClientSecret          string
	cfgTokenURL              string
	cfgScopes                string
	cfgAudience              string
	configureSecretsProvider secrets.SecretsProvider = secrets.DefaultSecretsProvider
)

var jwtConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure an OAuth client profile for JWT generation",
	Long:  `Add or update an OAuth client profile in your secrets for JWT generation.`,
	RunE:  runJWTConfigure,
}

func init() {
	jwtCmd.AddCommand(jwtConfigureCmd)

	jwtConfigureCmd.Flags().StringVarP(&cfgProfile, "profile", "p", "", "Profile name (required)")
	jwtConfigureCmd.Flags().StringVar(&cfgClientID, "client-id", "", "OAuth client ID")
	jwtConfigureCmd.Flags().StringVar(&cfgClientSecret, "client-secret", "", "OAuth client secret or private key")
	jwtConfigureCmd.Flags().StringVar(&cfgTokenURL, "token-url", "", "OAuth token URL")
	jwtConfigureCmd.Flags().StringVar(&cfgScopes, "scopes", "", "OAuth scopes (space-separated)")
	jwtConfigureCmd.Flags().StringVar(&cfgAudience, "audience", "", "OAuth audience")
}

func runJWTConfigure(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	if cfgProfile == "" {
		fmt.Print("Profile name: ")
		cfgProfile, _ = reader.ReadString('\n')
		cfgProfile = strings.TrimSpace(cfgProfile)
	}
	if cfgProfile == "" {
		return fmt.Errorf("profile name is required")
	}

	if cfgClientID == "" {
		fmt.Print("Client ID: ")
		cfgClientID, _ = reader.ReadString('\n')
		cfgClientID = strings.TrimSpace(cfgClientID)
	}
	if cfgClientSecret == "" {
		fmt.Print("Client Secret (or private key): ")
		cfgClientSecret, _ = reader.ReadString('\n')
		cfgClientSecret = strings.TrimSpace(cfgClientSecret)
	}
	if cfgTokenURL == "" {
		fmt.Print("Token URL: ")
		cfgTokenURL, _ = reader.ReadString('\n')
		cfgTokenURL = strings.TrimSpace(cfgTokenURL)
	}
	if cfgScopes == "" {
		fmt.Print("Scopes (space-separated): ")
		cfgScopes, _ = reader.ReadString('\n')
		cfgScopes = strings.TrimSpace(cfgScopes)
	}
	if cfgAudience == "" {
		fmt.Print("Audience: ")
		cfgAudience, _ = reader.ReadString('\n')
		cfgAudience = strings.TrimSpace(cfgAudience)
	}

	client := &secrets.OAuthClient{
		ClientID:     cfgClientID,
		ClientSecret: cfgClientSecret,
		TokenURL:     cfgTokenURL,
		Scopes:       cfgScopes,
		Audience:     cfgAudience,
	}

	err := configureSecretsProvider.StoreOAuthClient(cfgProfile, client)
	if err != nil {
		return fmt.Errorf("failed to store OAuth client: %w", err)
	}

	fmt.Printf("Profile '%s' configured successfully.\n", cfgProfile)
	return nil
}
