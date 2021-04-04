package n26api

import "os"

const (
	envUsername = "N26_USERNAME"
	envPassword = "N26_PASSWORD"
)

var _ CredentialsProvider = (*envCredentialsProvider)(nil)

// envCredentialsProvider provides username and password from environment variables.
type envCredentialsProvider struct{}

// Username provides a username from n26api.envUsername variable.
func (p *envCredentialsProvider) Username() string {
	return os.Getenv(envUsername)
}

// Password provides a password from n26api.envPassword variable.
func (p *envCredentialsProvider) Password() string {
	return os.Getenv(envPassword)
}

// CredentialsFromEnv initiates a new credentials provider.
func CredentialsFromEnv() CredentialsProvider {
	return &envCredentialsProvider{}
}
