package testkit

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.nhat.io/httpmock/test"
)

// TestingT is an alias of httpmock.TestingT.
type TestingT = test.T

// ServerOption is an option to configure Server.
type ServerOption = func(s *Server)

// ServerMocker is a function that applies expectations to the mocked server.
type ServerMocker func(t TestingT) *Server

// MockServer mocks a server with successful authentication workflow.
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
		s := NewServer(t)

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
