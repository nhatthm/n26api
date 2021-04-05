package testkit_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nhatthm/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/testkit"
	"github.com/nhatthm/n26api/pkg/util"
)

func TestMockServer(t *testing.T) {
	t.Parallel()

	authUsername := "nativeweb"
	authPassword := ""
	username := "username"
	password := "password"
	deviceID := uuid.New()

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", util.Base64Credentials(authUsername, authPassword)),
		"device-token":  deviceID.String(),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	s := testkit.MockServer(username, password, deviceID)(t)

	// 1st step: login.
	requestLoginBody := fmt.Sprintf("grant_type=password&password=%s&username=%s",
		url.QueryEscape(password), url.QueryEscape(username),
	)

	code, _, body, _ := request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestLoginBody))

	expectedLoginBody := fmt.Sprintf(
		`{"userMessage":{"title":"A second authentication factor is required.","detail":"Please provide your second form of authentication."},"mfaToken":%q,"error_description":"MFA token is required","detail":"MFA token is required","hostUrl":%q,"type":"mfa_required","error":"mfa_required","title":"A second authentication factor is required.","message":"Please provide your second form of authentication.","userId":%q,"status":403}`,
		s.MFAToken().String(),
		s.URL(),
		s.UserID(),
	)

	assert.NotEmpty(t, s.MFAToken())
	assert.Equal(t, http.StatusForbidden, code)
	assert.Equal(t, expectedLoginBody, string(body))

	// 2nd step: mfa challenge.
	requestMFAChallengeBody := fmt.Sprintf(`{"challengeType":"oob","mfaToken":%q}`, s.MFAToken().String())

	code, _, body, _ = request(t, s.URL(), http.MethodPost, "/api/mfa/challenge", requestHeader, []byte(requestMFAChallengeBody))

	expectedMFAChallengeBody := `{"challengeType":"oob"}`

	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, expectedMFAChallengeBody, string(body))

	// 3rd step: confirm login.
	requestConfirmLoginBody := fmt.Sprintf("grant_type=mfa_oob&mfaToken=%s",
		url.QueryEscape(s.MFAToken().String()),
	)

	code, _, body, _ = request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestConfirmLoginBody))

	expectedConfirmLoginBody := fmt.Sprintf(
		`{"access_token":%q,"token_type":"bearer","refresh_token":%q,"expires_in":889,"host_url":%q}`,
		s.AccessToken(),
		s.RefreshToken(),
		s.URL(),
	)

	assert.NotEmpty(t, s.AccessToken())
	assert.NotEmpty(t, s.RefreshToken())
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, expectedConfirmLoginBody, string(body))
}

func TestMockEmptyServer(t *testing.T) {
	t.Parallel()

	username := "nativeweb"
	password := ""
	mfaToken := uuid.New()
	accessToken := uuid.New()
	refreshToken := uuid.New()

	testCases := []struct {
		scenario      string
		mockServer    testkit.ServerMocker
		requestHeader map[string]string
		requestBody   string
	}{
		{
			scenario: "basic authentication",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthAuthorization(username, password),
				func(s *testkit.Server) {
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
			mockServer: testkit.MockEmptyServer(func(s *testkit.Server) {
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

// nolint:thelper
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
