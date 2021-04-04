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

// chain prepends the new provider to the chain.
func (chain *chainTokenProvider) chain(provider auth.TokenProvider) {
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

	for _, p := range providers {
		chain.chain(p)
	}

	return chain
}
