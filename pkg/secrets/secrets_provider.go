package secrets

// SecretsProvider defines the interface for secrets operations
// This allows for mocking in tests.
type SecretsProvider interface {
	GetOAuthClient(profile string) (*OAuthClient, error)
	StoreOAuthClient(profile string, client *OAuthClient) error
	ListOAuthProfiles() ([]string, error)
	DeleteOAuthClient(profile string) error
}

// DefaultSecretsProvider is the global secrets provider used in production
var DefaultSecretsProvider SecretsProvider = &RealSecrets{}

// For backward compatibility, keep the old function names as wrappers
func GetOAuthClient(profile string) (*OAuthClient, error) {
	return DefaultSecretsProvider.GetOAuthClient(profile)
}
func StoreOAuthClient(profile string, client *OAuthClient) error {
	return DefaultSecretsProvider.StoreOAuthClient(profile, client)
}
func ListOAuthProfiles() ([]string, error) {
	return DefaultSecretsProvider.ListOAuthProfiles()
}
func DeleteOAuthClient(profile string) error {
	return DefaultSecretsProvider.DeleteOAuthClient(profile)
}
