package testkit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/nhatthm/httpmock"
	"github.com/stretchr/testify/assert"

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

	s := MockEmptyServer(func(s *Server) {
		s.WithAuthAuthorization("nativeweb", "")
		s.ExpectWithBasicAuth(http.MethodGet, "/")
	})(t)

	expectedHeaders := httpmock.Header{
		"Authorization": "Basic bmF0aXZld2ViOg==",
	}

	assert.Equal(t, http.MethodGet, s.ExpectedRequests[0].Method)
	assert.Equal(t, httpmock.Exact("/"), s.ExpectedRequests[0].RequestURI)

	for key, matcher := range s.ExpectedRequests[0].RequestHeader {
		assert.True(t, matcher.Match(expectedHeaders[key]))
	}

	s.ResetExpectations()
}

func TestServer_Expect(t *testing.T) {
	t.Parallel()

	accessToken := uuid.New()

	s := MockEmptyServer(func(s *Server) {
		s.WithAccessToken(accessToken)
		s.Expect(http.MethodGet, "/")
	})(t)

	expectedHeaders := httpmock.Header{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	assert.Equal(t, http.MethodGet, s.ExpectedRequests[0].Method)
	assert.Equal(t, httpmock.Exact("/"), s.ExpectedRequests[0].RequestURI)

	for key, matcher := range s.ExpectedRequests[0].RequestHeader {
		assert.True(t, matcher.Match(expectedHeaders[key]))
	}

	s.ResetExpectations()
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

			s := MockEmptyServer(tc.mockServer)(t)
			s.WithAccessToken(accessToken)

			expectedHeaders := httpmock.Header{
				"Authorization": fmt.Sprintf("Bearer %s", accessToken),
			}

			assert.Equal(t, tc.expectedMethod, s.ExpectedRequests[0].Method)
			assert.Equal(t, httpmock.Exact("/"), s.ExpectedRequests[0].RequestURI)

			for key, matcher := range s.ExpectedRequests[0].RequestHeader {
				assert.True(t, matcher.Match(expectedHeaders[key]))
			}

			s.ResetExpectations()
		})
	}
}
