package n26api

// CredentialsProvider provides username and password for authentication.
type CredentialsProvider interface {
	// Username provides a username.
	Username() string
	// Password provides a password.
	Password() string
}
