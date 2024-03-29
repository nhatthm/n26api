package n26api_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.nhat.io/httpmock"

	"github.com/nhatthm/n26api"
	"github.com/nhatthm/n26api/pkg/testkit/auth"
)

func TestRoundTripper(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		tripper        func(*testing.T) n26api.RoundTripperFunc
		expectedHeader string
	}{
		{
			scenario: "basic",
			tripper: func(*testing.T) n26api.RoundTripperFunc {
				return n26api.BasicAuthRoundTripper("foo", "bar", http.DefaultTransport)
			},
			expectedHeader: "Basic Zm9vOmJhcg==",
		},
		{
			scenario: "bearer",
			tripper: func(*testing.T) n26api.RoundTripperFunc {
				return n26api.BearerAuthRoundTripper("foobar", http.DefaultTransport)
			},
			expectedHeader: "Bearer foobar",
		},
		{
			scenario: "token",
			tripper: func(*testing.T) n26api.RoundTripperFunc {
				tokenProvider := auth.MockTokenProvider(func(p *auth.TokenProvider) {
					p.On("Token", mock.Anything).
						Return("foobaz", nil)
				})(t)

				return n26api.TokenRoundTripper(tokenProvider, http.DefaultTransport)
			},
			expectedHeader: "Bearer foobaz",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/").
					WithHeader("Authorization", tc.expectedHeader).
					Return("hello world!")
			})(t)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.URL(), nil)
			require.NoError(t, err, "could not create a new request")

			client := http.Client{
				Timeout:   time.Second,
				Transport: tc.tripper(t),
			}
			resp, err := client.Do(req)
			require.NoError(t, err, "could not make a request to mocked server")

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "could not read response body")

			err = resp.Body.Close()
			require.NoError(t, err, "could not close response body")

			expectedBody := `hello world!`

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, expectedBody, string(respBody))
		})
	}
}

func TestRoundTripper_Error(t *testing.T) {
	t.Parallel()

	p := auth.MockTokenProvider(func(p *auth.TokenProvider) {
		p.On("Token", mock.Anything).
			Return("", errors.New("token error"))
	})(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := n26api.TokenRoundTripper(p, nil)(req) // nolint:bodyclose

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "token error")
}
