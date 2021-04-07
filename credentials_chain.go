package n26api

var _ CredentialsProvider = (*chainCredentialsProvider)(nil)

type chainCredentialsProvider []CredentialsProvider

// Username provides the first non-empty username from the chain.
func (chain *chainCredentialsProvider) Username() string {
	for _, p := range *chain {
		if username := p.Username(); username != "" {
			return username
		}
	}

	return ""
}

// Password provides the first non-empty password from the chain.
func (chain *chainCredentialsProvider) Password() string {
	for _, p := range *chain {
		if username := p.Password(); username != "" {
			return username
		}
	}

	return ""
}

// append appends a new provider to the chain.
func (chain *chainCredentialsProvider) append(provider CredentialsProvider) {
	*chain = append(*chain, provider)
}

// prepend prepends a new provider to the chain.
func (chain *chainCredentialsProvider) prepend(provider CredentialsProvider) {
	*chain = append(*chain, provider)
	copy((*chain)[1:], *chain)
	(*chain)[0] = provider
}

// newChainCredentialsProvider initiates a chain of CredentialsProvider.
func newChainCredentialsProvider() *chainCredentialsProvider {
	chain := make(chainCredentialsProvider, 0)

	return &chain
}

// chainCredentialsProviders chains a list of CredentialsProvider.
func chainCredentialsProviders(providers ...CredentialsProvider) *chainCredentialsProvider {
	chain := newChainCredentialsProvider()
	*chain = providers

	return chain
}
