package testkit_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/testkit"
	"github.com/nhatthm/n26api/pkg/util"
)

func TestWithAuthPasswordLogin_Error(t *testing.T) {
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
	requestBody := fmt.Sprintf("grant_type=password&password=%s&username=%s",
		url.QueryEscape(password), url.QueryEscape(username),
	)

	testCases := []struct {
		scenario     string
		option       testkit.ServerOption
		expectedCode int
		expectedBody string
	}{
		{
			scenario:     "wrong credentials",
			option:       testkit.WithAuthPasswordLoginFailureWrongCredentials(username, password, deviceID),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":400,"detail":"Bad credentials","type":"invalid_grant","userMessage":{"title":"Login failed!","detail":"Error! The email address or password is incorrect"},"error":"invalid_grant","error_description":"Bad credentials"}`,
		},
		{
			scenario:     "too many attempts",
			option:       testkit.WithAuthPasswordLoginFailureTooManyAttempts(username, password, deviceID),
			expectedCode: http.StatusTooManyRequests,
			expectedBody: `{"status":429,"detail":"Too many log-in attempts. Please try again in 30 minutes.","userMessage":{"title":"Oops!","detail":"Too many log-in attempts. Please try again in 30 minutes."},"error":"Oops!","title":"Oops!","message":"Too many log-in attempts. Please try again in 30 minutes."}`,
		},
		{
			scenario:     "unexpected response",
			option:       testkit.WithAuthPasswordLoginUnexpectedResponse(username, password, deviceID),
			expectedCode: http.StatusOK,
		},
		{
			scenario:     "internal server error",
			option:       testkit.WithAuthPasswordLoginFailure(username, password, deviceID),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := testkit.MockEmptyServer(
				testkit.WithAuthAuthorization(authUsername, authPassword),
				tc.option,
			)(t)

			code, _, body, _ := request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestBody))

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}

func TestWithAuthPasswordLoginSuccess(t *testing.T) {
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
	requestBody := fmt.Sprintf("grant_type=password&password=%s&username=%s",
		url.QueryEscape(password), url.QueryEscape(username),
	)

	s := testkit.MockEmptyServer(
		testkit.WithAuthAuthorization(authUsername, authPassword),
		testkit.WithAuthPasswordLoginSuccess(username, password, deviceID),
	)(t)

	code, _, body, _ := request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestBody))

	expectedBody := fmt.Sprintf(
		`{"userMessage":{"title":"A second authentication factor is required.","detail":"Please provide your second form of authentication."},"mfaToken":%q,"error_description":"MFA token is required","detail":"MFA token is required","hostUrl":%q,"type":"mfa_required","error":"mfa_required","title":"A second authentication factor is required.","message":"Please provide your second form of authentication.","userId":%q,"status":403}`,
		s.MFAToken().String(),
		s.URL(),
		s.UserID(),
	)

	assert.NotEmpty(t, s.MFAToken())
	assert.Equal(t, http.StatusForbidden, code)
	assert.Equal(t, expectedBody, string(body))
}

func TestWithAuthMFAChallenge_Error(t *testing.T) {
	t.Parallel()

	authUsername := "nativeweb"
	authPassword := ""
	deviceID := uuid.New()
	mfaToken := uuid.New()

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", util.Base64Credentials(authUsername, authPassword)),
		"device-token":  deviceID.String(),
	}
	requestBody := fmt.Sprintf(`{"challengeType":"oob","mfaToken":%q}`, mfaToken.String())

	testCases := []struct {
		scenario     string
		option       testkit.ServerOption
		expectedCode int
		expectedBody string
	}{
		{
			scenario:     "invalid token",
			option:       testkit.WithAuthMFAChallengeFailureInvalidToken(),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":401,"detail":"Invalid token","type":"error","userMessage":{"title":"Login attempt expired","detail":"That took too long, please try again."},"error":"invalid_token","error_description":"Invalid token"}`,
		},
		{
			scenario:     "internal server error",
			option:       testkit.WithAuthMFAChallengeFailure(),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := testkit.MockEmptyServer(
				testkit.WithAuthAuthorization(authUsername, authPassword),
				func(s *testkit.Server) {
					s.WithDeviceID(deviceID).
						WithMFAToken(mfaToken)
				},
				tc.option,
			)(t)

			code, _, body, _ := request(t, s.URL(), http.MethodPost, "/api/mfa/challenge", requestHeader, []byte(requestBody))

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}

func TestWithAuthMFAChallengeSuccess(t *testing.T) {
	t.Parallel()

	authUsername := "nativeweb"
	authPassword := ""
	deviceID := uuid.New()
	mfaToken := uuid.New()

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", util.Base64Credentials(authUsername, authPassword)),
		"device-token":  deviceID.String(),
	}
	requestBody := fmt.Sprintf(`{"challengeType":"oob","mfaToken":%q}`, mfaToken.String())

	s := testkit.MockEmptyServer(
		testkit.WithAuthAuthorization(authUsername, authPassword),
		func(s *testkit.Server) {
			s.WithDeviceID(deviceID).
				WithMFAToken(mfaToken)
		},
		testkit.WithAuthMFAChallengeSuccess(),
	)(t)

	code, _, body, _ := request(t, s.URL(), http.MethodPost, "/api/mfa/challenge", requestHeader, []byte(requestBody))

	expectedBody := `{"challengeType":"oob"}`

	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, expectedBody, string(body))
}

func TestWithAuthConfirmLogin_Error(t *testing.T) {
	t.Parallel()

	authUsername := "nativeweb"
	authPassword := ""
	deviceID := uuid.New()
	mfaToken := uuid.New()

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", util.Base64Credentials(authUsername, authPassword)),
		"device-token":  deviceID.String(),
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	requestBody := fmt.Sprintf("grant_type=mfa_oob&mfaToken=%s",
		url.QueryEscape(mfaToken.String()),
	)

	testCases := []struct {
		scenario     string
		option       testkit.ServerOption
		expectedCode int
		expectedBody string
	}{
		{
			scenario:     "invalid token",
			option:       testkit.WithAuthConfirmLoginFailureInvalidToken(1),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":401,"detail":"Invalid token","type":"error","userMessage":{"title":"Login attempt expired","detail":"That took too long, please try again."},"error":"invalid_token","error_description":"Invalid token"}`,
		},
		{
			scenario:     "internal server error",
			option:       testkit.WithAuthConfirmLoginFailure(),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := testkit.MockEmptyServer(
				testkit.WithAuthAuthorization(authUsername, authPassword),
				func(s *testkit.Server) {
					s.WithDeviceID(deviceID).
						WithMFAToken(mfaToken)
				},
				tc.option,
			)(t)

			code, _, body, _ := request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestBody))

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}

func TestWithAuthConfirmLoginSuccess(t *testing.T) {
	t.Parallel()

	authUsername := "nativeweb"
	authPassword := ""
	deviceID := uuid.New()
	mfaToken := uuid.New()

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", util.Base64Credentials(authUsername, authPassword)),
		"device-token":  deviceID.String(),
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	requestBody := fmt.Sprintf("grant_type=mfa_oob&mfaToken=%s",
		url.QueryEscape(mfaToken.String()),
	)

	s := testkit.MockEmptyServer(
		testkit.WithAuthAuthorization(authUsername, authPassword),
		func(s *testkit.Server) {
			s.WithDeviceID(deviceID).
				WithMFAToken(mfaToken)
		},
		testkit.WithAuthConfirmLoginSuccess(),
	)(t)

	code, _, body, _ := request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestBody))

	expectedBody := fmt.Sprintf(
		`{"access_token":%q,"token_type":"bearer","refresh_token":%q,"expires_in":889,"host_url":%q}`,
		s.AccessToken(),
		s.RefreshToken(),
		s.URL(),
	)

	assert.NotEmpty(t, s.AccessToken())
	assert.NotEmpty(t, s.RefreshToken())
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, expectedBody, string(body))
}

func TestWithAuthRefreshToken_Error(t *testing.T) {
	t.Parallel()

	authUsername := "nativeweb"
	authPassword := ""
	deviceID := uuid.New()
	refreshToken := uuid.New()

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", util.Base64Credentials(authUsername, authPassword)),
		"device-token":  deviceID.String(),
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	requestBody := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s",
		url.QueryEscape(refreshToken.String()),
	)

	testCases := []struct {
		scenario     string
		option       testkit.ServerOption
		expectedCode int
		expectedBody string
	}{
		{
			scenario:     "invalid token",
			option:       testkit.WithAuthRefreshTokenFailureInvalidToken(),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"status":401,"detail":"Invalid token","type":"error","userMessage":{"title":"Login attempt expired","detail":"That took too long, please try again."},"error":"invalid_token","error_description":"Invalid token"}`,
		},
		{
			scenario:     "internal server error",
			option:       testkit.WithAuthRefreshTokenFailure(),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := testkit.MockEmptyServer(
				testkit.WithAuthAuthorization(authUsername, authPassword),
				func(s *testkit.Server) {
					s.WithDeviceID(deviceID).
						WithRefreshToken(refreshToken)
				},
				tc.option,
			)(t)

			code, _, body, _ := request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestBody))

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}

func TestWithAuthRefreshTokenSuccess(t *testing.T) {
	t.Parallel()

	authUsername := "nativeweb"
	authPassword := ""
	deviceID := uuid.New()
	refreshToken := uuid.New()

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", util.Base64Credentials(authUsername, authPassword)),
		"device-token":  deviceID.String(),
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	requestBody := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s",
		url.QueryEscape(refreshToken.String()),
	)

	s := testkit.MockEmptyServer(
		testkit.WithAuthAuthorization(authUsername, authPassword),
		func(s *testkit.Server) {
			s.WithDeviceID(deviceID).
				WithRefreshToken(refreshToken)
		},
		testkit.WithAuthRefreshTokenSuccess(),
	)(t)

	code, _, body, _ := request(t, s.URL(), http.MethodPost, "/oauth/token", requestHeader, []byte(requestBody))

	expectedBody := fmt.Sprintf(
		`{"access_token":%q,"token_type":"bearer","refresh_token":%q,"expires_in":889,"host_url":%q}`,
		s.AccessToken(),
		s.RefreshToken(),
		s.URL(),
	)

	assert.NotEmpty(t, s.AccessToken())
	assert.NotEmpty(t, s.RefreshToken())
	assert.NotEqual(t, refreshToken.String(), string(s.RefreshToken()))
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, expectedBody, string(body))
}
