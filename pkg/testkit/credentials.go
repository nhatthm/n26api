package testkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CredentialsProviderMocker is CredentialsProvider mocker.
type CredentialsProviderMocker func(tb testing.TB) *CredentialsProvider

// NoMockCredentialsProvider is no mock CredentialsProvider.
var NoMockCredentialsProvider = MockCredentialsProvider()

// CredentialsProvider is a CredentialsProvider.
type CredentialsProvider struct {
	mock.Mock
}

// Username satisfies CredentialsProvider.
func (c *CredentialsProvider) Username() string {
	return c.Called().String(0)
}

// Password satisfies CredentialsProvider.
func (c *CredentialsProvider) Password() string {
	return c.Called().String(0)
}

// mockCredentialsProvider mocks CredentialsProvider interface.
func mockCredentialsProvider(mocks ...func(p *CredentialsProvider)) *CredentialsProvider {
	p := &CredentialsProvider{}

	for _, m := range mocks {
		m(p)
	}

	return p
}

// MockCredentialsProvider creates CredentialsProvider mock with cleanup to ensure all the expectations are met.
func MockCredentialsProvider(mocks ...func(p *CredentialsProvider)) CredentialsProviderMocker {
	return func(tb testing.TB) *CredentialsProvider {
		tb.Helper()

		p := mockCredentialsProvider(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, p.Mock.AssertExpectations(tb))
		})

		return p
	}
}
