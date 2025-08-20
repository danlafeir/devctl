package secrets

// SecretsProvider defines the generic interface for secrets operations
// This allows for mocking in tests and provides standard CRUD operations.
// Keys are stored with the naming convention: cli.devctl.<cmd>.<token>
type SecretsProvider interface {
	Read(cmd, token string) (string, error)
	Write(cmd, token, value string) error
	List(cmd string) ([]string, error)
	Delete(cmd, token string) error
}

// DefaultSecretsProvider is the global secrets provider used in production
var DefaultSecretsProvider SecretsProvider = &RealSecrets{}
