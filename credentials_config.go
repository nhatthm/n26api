package n26api

var _ CredentialsProvider = (*configCredentialsProvider)(nil)

// configCredentialsProvider provides username and password from config.
type configCredentialsProvider struct {
	username string
	password string
}

// Username provides a username from config.
func (p *configCredentialsProvider) Username() string {
	return p.username
}

// Password provides a password from config.
func (p *configCredentialsProvider) Password() string {
	return p.password
}

// Credentials initiates a new credentials provider.
func Credentials(username string, password string) CredentialsProvider {
	return &configCredentialsProvider{
		username: username,
		password: password,
	}
}
