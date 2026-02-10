package secrets

// SecretsProvider defines the generic interface for secrets operations
// This allows for mocking in tests and provides standard CRUD operations.
// Keys are stored with the naming convention: cli.devctl.<namespace>.<name>
type SecretsProvider interface {
	Read(namespace, name string) (string, error)
	Write(namespace, name, value string) error
	List(namespace string) ([]string, error)
	Delete(namespace, name string) error
}

// defaultProvider is the secrets provider used by package-level functions
var defaultProvider SecretsProvider = &RealSecrets{}

// Read retrieves a secret value by namespace and name
func Read(namespace, name string) (string, error) {
	return defaultProvider.Read(namespace, name)
}

// Write stores a secret value by namespace and name
func Write(namespace, name, value string) error {
	return defaultProvider.Write(namespace, name, value)
}

// List returns all secret names within a namespace
func List(namespace string) ([]string, error) {
	return defaultProvider.List(namespace)
}

// Delete removes a secret by namespace and name
func Delete(namespace, name string) error {
	return defaultProvider.Delete(namespace, name)
}
