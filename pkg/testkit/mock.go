package testkit

import (
	"strings"

	"github.com/google/uuid"
	"github.com/nhatthm/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/auth"
)

// TestingT is an alias of httpmock.TestingT.
type TestingT = httpmock.TestingT

// ServerOption is an option to configure Server.
type ServerOption = func(s *Server)

// ServerMocker is a function that applies expectations to the mocked server.
type ServerMocker func(t TestingT) *Server

// MockServer mocks a server with authentication.
func MockServer(
	username, password string,
	deviceID uuid.UUID,
	mocks ...ServerOption,
) ServerMocker {
	defaults := []ServerOption{
		WithAuthSuccess(username, password, deviceID),
	}

	args := make([]ServerOption, 0, len(mocks)+len(defaults))
	args = append(args, defaults...)
	args = append(args, mocks...)

	return MockEmptyServer(args...)
}

// MockEmptyServer mocks a N26 API server.
func MockEmptyServer(mocks ...ServerOption) ServerMocker {
	return func(t TestingT) *Server {
		s := &Server{
			Server: httpmock.NewServer(t),
			userID: uuid.New(),
		}

		s.WithAuthAuthorization(auth.BasicAuthUsername, auth.BasicAuthPassword).
			WithDefaultResponseHeaders(httpmock.Header{
				"Content-Type": "application/json",
			})

		s.WithRequestMatcher(
			httpmock.SequentialRequestMatcher(
				httpmock.WithBodyMatcher(func(t TestingT, expected, body []byte) bool {
					replaced := strings.ReplaceAll(string(expected), "{{MFAToken}}", s.mfaToken.String())
					replaced = strings.ReplaceAll(replaced, "{{RefreshToken}}", string(s.refreshToken))

					return assert.Equal(t, []byte(replaced), body)
				}),
				httpmock.WithHeaderMatcher(func(t httpmock.TestingT, expected, header string) bool {
					replaced := strings.ReplaceAll(expected, "{{accessToken}}", string(s.accessToken))

					return assert.Equal(t, replaced, header)
				}),
			),
		)

		for _, m := range mocks {
			m(s)
		}

		t.Cleanup(func() {
			assert.NoError(t, s.ExpectationsWereMet())
			s.Close()
		})

		return s
	}
}

// WithAuthAuthorization sets the Authorization credentials for asserting the /oauth/token request.
func WithAuthAuthorization(username, password string) ServerOption {
	return func(s *Server) {
		s.WithAuthAuthorization(username, password)
	}
}
