package testkit

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nhatthm/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMockServer(t *testing.T) {
	t.Parallel()

	username := "nativeweb"
	password := ""
	mfaToken := uuid.New()
	accessToken := uuid.New()
	refreshToken := uuid.New()

	testCases := []struct {
		scenario      string
		mockServer    ServerMocker
		requestHeader map[string]string
		requestBody   string
	}{
		{
			scenario: "basic authentication",
			mockServer: MockEmptyServer(
				WithAuthAuthorization(username, password),
				func(s *Server) {
					s.WithMFAToken(mfaToken)
					s.WithRefreshToken(refreshToken)
					s.ExpectWithBasicAuth(http.MethodGet, "/").
						WithBody("MFA Token: {{MFAToken}}\nRefresh Token: {{RefreshToken}}").
						Return(`{}`)
				},
			),
			requestHeader: map[string]string{"Authorization": "Basic bmF0aXZld2ViOg=="},
			requestBody:   fmt.Sprintf("MFA Token: %s\nRefresh Token: %s", mfaToken.String(), refreshToken.String()),
		},
		{
			scenario: "bearer authentication",
			mockServer: MockEmptyServer(func(s *Server) {
				s.WithAccessToken(accessToken)
				s.ExpectGet("/").Return(`{}`)
			}),
			requestHeader: map[string]string{"Authorization": fmt.Sprintf(`Bearer %s`, accessToken.String())},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockServer(t)
			code, headers, _, _ := request(t, s.URL(), http.MethodGet, "/", tc.requestHeader, []byte(tc.requestBody))

			expectedHeaders := map[string]string{
				"Content-Type": "application/json",
			}

			assert.Equal(t, http.StatusOK, code)
			httpmock.AssertHeaderContains(t, headers, expectedHeaders)
		})
	}
}

func request(
	t *testing.T,
	baseURL string,
	method, uri string,
	headers map[string]string,
	body []byte,
) (int, map[string]string, []byte, time.Duration) {
	return httpmock.DoRequest(t,
		method, baseURL+uri,
		headers, body,
	)
}
