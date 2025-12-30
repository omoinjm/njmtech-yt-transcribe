package transcriber

import (
	"os"
)

// EnvAPIKeyProvider implements the APIKeyProvider interface by
// retrieving the API key from an environment variable.
// This decouples the transcriber from the specific method of fetching credentials.
type EnvAPIKeyProvider struct {
	EnvVarName string // The name of the environment variable to read
}

// NewEnvAPIKeyProvider creates and returns a new instance of EnvAPIKeyProvider.
func NewEnvAPIKeyProvider(envVarName string) *EnvAPIKeyProvider {
	return &EnvAPIKeyProvider{
		EnvVarName: envVarName,
	}
}

// GetAPIKey retrieves the API key from the environment variable specified by EnvVarName.
func (p *EnvAPIKeyProvider) GetAPIKey() string {
	return os.Getenv(p.EnvVarName)
}
