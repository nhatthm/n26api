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

// chain prepends the new provider to the chain.
func (chain *chainCredentialsProvider) chain(provider CredentialsProvider) {
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

	for _, p := range providers {
		chain.chain(p)
	}

	return chain
}
