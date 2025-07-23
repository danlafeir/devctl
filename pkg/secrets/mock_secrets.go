package secrets

import "fmt"

// MockSecrets is an in-memory implementation for tests
// Not thread-safe, but sufficient for unit tests
type MockSecrets struct {
	store map[string]*OAuthClient
}

func NewMockSecrets() *MockSecrets {
	return &MockSecrets{store: make(map[string]*OAuthClient)}
}

func (m *MockSecrets) GetOAuthClient(profile string) (*OAuthClient, error) {
	c, ok := m.store[profile]
	if !ok {
		return nil, fmt.Errorf("profile not found")
	}
	return c, nil
}
func (m *MockSecrets) StoreOAuthClient(profile string, client *OAuthClient) error {
	m.store[profile] = client
	return nil
}
func (m *MockSecrets) ListOAuthProfiles() ([]string, error) {
	profiles := make([]string, 0, len(m.store))
	for k := range m.store {
		profiles = append(profiles, k)
	}
	return profiles, nil
}
func (m *MockSecrets) DeleteOAuthClient(profile string) error {
	delete(m.store, profile)
	return nil
}
