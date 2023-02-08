package testkit

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.nhat.io/httpmock"
	httpMock "go.nhat.io/httpmock/mock/http"
	plannerMock "go.nhat.io/httpmock/mock/planner"
	"go.nhat.io/httpmock/planner"

	"github.com/nhatthm/n26api/pkg/auth"
)

func TestServer_WithAuthAuthorization(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		username string
		password string
		expected string
	}{
		{
			scenario: "empty",
			expected: "Basic",
		},
		{
			scenario: "password is empty",
			username: "nativeweb",
			expected: "Basic bmF0aXZld2ViOg==",
		},
		{
			scenario: "not empty",
			username: "username",
			password: "password",
			expected: "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := (&Server{}).WithAuthAuthorization(tc.username, tc.password)

			assert.Equal(t, tc.expected, s.BasicAuthorization())
		})
	}
}

func TestServer_UserID(t *testing.T) {
	t.Parallel()

	expected := uuid.New()

	s := &Server{userID: expected}

	assert.Equal(t, expected, s.UserID())
}

func TestServer_WithDeviceID(t *testing.T) {
	t.Parallel()

	expected := uuid.New()

	s := (&Server{}).WithDeviceID(expected)

	assert.Equal(t, expected, s.DeviceID())
}

func TestServer_WithMFAToken(t *testing.T) {
	t.Parallel()

	expected := uuid.New()

	s := (&Server{}).WithMFAToken(expected)

	assert.Equal(t, expected, s.MFAToken())
}

func TestServer_WithAccessToken(t *testing.T) {
	t.Parallel()

	token := uuid.New()
	expected := auth.Token(token.String())

	s := (&Server{}).WithAccessToken(token)

	assert.Equal(t, expected, s.AccessToken())
}

func TestServer_WithRefreshToken(t *testing.T) {
	t.Parallel()

	token := uuid.New()
	expected := auth.Token(token.String())

	s := (&Server{}).WithRefreshToken(token)

	assert.Equal(t, expected, s.RefreshToken())
}

func TestServer_ExpectWithBasicAuth(t *testing.T) {
	t.Parallel()

	var e planner.Expectation

	p := plannerMock.Mock(func(p *plannerMock.Planner) {
		p.On("Expect", mock.Anything).
			Run(func(args mock.Arguments) {
				e = args[0].(planner.Expectation)
			})

		p.On("IsEmpty").Return(true)
	})(t)

	MockEmptyServer(func(s *Server) {
		s.WithPlanner(p).
			WithAuthAuthorization("nativeweb", "")

		s.ExpectWithBasicAuth(http.MethodGet, "/")
	})(t)

	expectedHeaders := httpmock.Header{
		"Authorization": "Basic bmF0aXZld2ViOg==",
	}

	assert.Equal(t, http.MethodGet, e.Method())
	assert.Equal(t, httpmock.Exact("/"), e.URIMatcher())

	requestHeader := e.HeaderMatcher()

	assert.Len(t, requestHeader, 1)

	for key, m := range requestHeader {
		matched, err := m.Match(expectedHeaders[key])

		assert.True(t, matched)
		assert.NoError(t, err)
	}
}

func TestServer_Expect(t *testing.T) {
	t.Parallel()

	accessToken := uuid.New()

	var e planner.Expectation

	p := plannerMock.Mock(func(p *plannerMock.Planner) {
		p.On("Expect", mock.Anything).
			Run(func(args mock.Arguments) {
				e = args[0].(planner.Expectation)
			})

		p.On("IsEmpty").Return(true)
	})(t)

	MockEmptyServer(func(s *Server) {
		s.WithPlanner(p).
			WithAccessToken(accessToken)

		s.Expect(http.MethodGet, "/")
	})(t)

	expectedHeaders := httpmock.Header{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	assert.Equal(t, http.MethodGet, e.Method())
	assert.Equal(t, httpmock.Exact("/"), e.URIMatcher())

	requestHeader := e.HeaderMatcher()

	assert.Len(t, requestHeader, 1)

	for key, m := range requestHeader {
		matched, err := m.Match(expectedHeaders[key])

		assert.True(t, matched)
		assert.NoError(t, err)
	}
}

func TestServer_ExpectAliases(t *testing.T) {
	t.Parallel()

	accessToken := uuid.New()

	testCases := []struct {
		scenario       string
		mockServer     func(s *Server)
		expectedMethod string
	}{
		{
			scenario: "GET",
			mockServer: func(s *Server) {
				s.ExpectGet("/")
			},
			expectedMethod: http.MethodGet,
		},
		{
			scenario: "HEAD",
			mockServer: func(s *Server) {
				s.ExpectHead("/")
			},
			expectedMethod: http.MethodHead,
		},
		{
			scenario: "POST",
			mockServer: func(s *Server) {
				s.ExpectPost("/")
			},
			expectedMethod: http.MethodPost,
		},
		{
			scenario: "PUT",
			mockServer: func(s *Server) {
				s.ExpectPut("/")
			},
			expectedMethod: http.MethodPut,
		},
		{
			scenario: "PATCH",
			mockServer: func(s *Server) {
				s.ExpectPatch("/")
			},
			expectedMethod: http.MethodPatch,
		},
		{
			scenario: "DELETE",
			mockServer: func(s *Server) {
				s.ExpectDelete("/")
			},
			expectedMethod: http.MethodDelete,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var e planner.Expectation

			p := plannerMock.Mock(func(p *plannerMock.Planner) {
				p.On("Expect", mock.Anything).
					Run(func(args mock.Arguments) {
						e = args[0].(planner.Expectation)
					})

				p.On("IsEmpty").Return(true)
			})(t)

			s := MockEmptyServer(func(s *Server) {
				s.WithPlanner(p)
			}, tc.mockServer)(t)

			s.WithAccessToken(accessToken)

			expectedHeaders := httpmock.Header{
				"Authorization": fmt.Sprintf("Bearer %s", accessToken),
			}

			assert.Equal(t, tc.expectedMethod, e.Method())
			assert.Equal(t, httpmock.Exact("/"), e.URIMatcher())

			requestHeader := e.HeaderMatcher()

			assert.Len(t, requestHeader, 1)

			for key, m := range requestHeader {
				matched, err := m.Match(expectedHeaders[key])

				assert.True(t, matched)
				assert.NoError(t, err)
			}
		})
	}
}

func handleRequestSuccess(t *testing.T, h httpmock.ExpectationHandler) ([]byte, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	req, _ := http.NewRequest(http.MethodGet, "/", nil) //nolint: errcheck

	w := httpMock.MockResponseWriter(func(w *httpMock.ResponseWriter) {
		w.On("WriteHeader", httpmock.StatusOK)

		w.On("Write", mock.Anything).
			Run(func(args mock.Arguments) {
				buf.Write(args[0].([]byte))
			}).
			Return(0, nil)
	})(t)

	err := h.Handle(w, req, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
