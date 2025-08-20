package mocks

import (
	"fmt"
	"strings"
)

// MockSecrets is an in-memory implementation for tests
// Not thread-safe, but sufficient for unit tests
type MockSecrets struct {
	store map[string]string
}

func NewMockSecrets() *MockSecrets {
	return &MockSecrets{store: make(map[string]string)}
}

func (m *MockSecrets) Read(cmd, token string) (string, error) {
	key := fmt.Sprintf("cli.devctl.%s.%s", cmd, token)
	value, ok := m.store[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

func (m *MockSecrets) Write(cmd, token, value string) error {
	key := fmt.Sprintf("cli.devctl.%s.%s", cmd, token)
	m.store[key] = value
	return nil
}

func (m *MockSecrets) List(cmd string) ([]string, error) {
	prefix := fmt.Sprintf("cli.devctl.%s.", cmd)
	var tokens []string
	for k := range m.store {
		if strings.HasPrefix(k, prefix) {
			token := strings.TrimPrefix(k, prefix)
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}

func (m *MockSecrets) Delete(cmd, token string) error {
	key := fmt.Sprintf("cli.devctl.%s.%s", cmd, token)
	delete(m.store, key)
	return nil
}
