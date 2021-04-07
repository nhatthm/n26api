package n26api

import (
	"context"

	"github.com/nhatthm/n26api/pkg/auth"
)

var _ auth.TokenProvider = (*chainTokenProvider)(nil)

type chainTokenProvider []auth.TokenProvider

// Token provides the first non-empty token from the chain.
func (chain *chainTokenProvider) Token(ctx context.Context) (auth.Token, error) {
	for _, p := range *chain {
		token, err := p.Token(ctx)
		if err != nil {
			return "", err
		}

		if token != "" {
			return token, nil
		}
	}

	return "", nil
}

// append appends a new provider to the chain.
func (chain *chainTokenProvider) append(provider auth.TokenProvider) {
	*chain = append(*chain, provider)
}

// prepend prepends a new provider to the chain.
func (chain *chainTokenProvider) prepend(provider auth.TokenProvider) {
	*chain = append(*chain, provider)
	copy((*chain)[1:], *chain)
	(*chain)[0] = provider
}

// newChainTokenProvider initiates a chain of auth.TokenProvider.
func newChainTokenProvider() *chainTokenProvider {
	chain := make(chainTokenProvider, 0)

	return &chain
}

// chainTokenProviders chains a list of auth.TokenProvider.
func chainTokenProviders(providers ...auth.TokenProvider) *chainTokenProvider {
	chain := newChainTokenProvider()
	*chain = providers

	return chain
}
