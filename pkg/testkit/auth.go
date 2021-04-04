package testkit

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/nhatthm/n26api/internal/api"
)

func expectAuthPasswordLogin(s *Server, username, password string, deviceID uuid.UUID) *Request {
	return s.WithDeviceID(deviceID).
		ExpectWithBasicAuth(http.MethodPost, "/oauth/token").
		WithHeader("device-token", s.DeviceID().String()).
		WithHeader("Content-Type", "application/x-www-form-urlencoded").
		WithBodyf("grant_type=password&password=%s&username=%s",
			url.QueryEscape(password), url.QueryEscape(username),
		)
}

func expectMFAChallenge(s *Server) *Request {
	return s.ExpectWithBasicAuth(http.MethodPost, "/api/mfa/challenge").
		WithHeader("device-token", s.DeviceID().String()).
		WithBody(`{"challengeType":"oob","mfaToken":"{{MFAToken}}"}`)
}

func expectConfirmLogin(s *Server) *Request {
	return s.ExpectWithBasicAuth(http.MethodPost, "/oauth/token").
		WithHeader("device-token", s.DeviceID().String()).
		WithHeader("Content-Type", "application/x-www-form-urlencoded").
		WithBody("grant_type=mfa_oob&mfaToken={{MFAToken}}")
}

func expectRefreshToken(s *Server) *Request {
	return s.ExpectWithBasicAuth(http.MethodPost, "/oauth/token").
		WithHeader("device-token", s.DeviceID().String()).
		WithHeader("Content-Type", "application/x-www-form-urlencoded").
		WithBody("grant_type=refresh_token&refresh_token={{RefreshToken}}")
}

// WithAuthPasswordLoginFailureWrongCredentials expects a request for login and returns a bad credentials error (400).
func WithAuthPasswordLoginFailureWrongCredentials(username, password string, deviceID uuid.UUID) ServerOption {
	return func(s *Server) {
		expectAuthPasswordLogin(s, username, password, deviceID).
			ReturnCode(http.StatusBadRequest).
			ReturnJSON(api.BadCredentialsError{
				Status: http.StatusBadRequest,
				Detail: "Bad credentials",
				Type:   "invalid_grant",
				UserMessage: api.UserMessage{
					Title:  "Login failed!",
					Detail: "Error! The email address or password is incorrect",
				},
				Error:            "invalid_grant",
				ErrorDescription: "Bad credentials",
			})
	}
}

// WithAuthPasswordLoginFailureTooManyAttempts expects a request for login and returns a too many attempts error (429).
func WithAuthPasswordLoginFailureTooManyAttempts(username, password string, deviceID uuid.UUID) ServerOption {
	return func(s *Server) {
		expectAuthPasswordLogin(s, username, password, deviceID).
			ReturnCode(http.StatusTooManyRequests).
			ReturnJSON(api.TooManyLoginAttemptsError{
				Status: http.StatusTooManyRequests,
				Detail: "Too many log-in attempts. Please try again in 30 minutes.",
				UserMessage: api.UserMessage{
					Title:  "Oops!",
					Detail: "Too many log-in attempts. Please try again in 30 minutes.",
				},
				Error:   "Oops!",
				Title:   "Oops!",
				Message: "Too many log-in attempts. Please try again in 30 minutes.",
			})
	}
}

// WithAuthPasswordLoginUnexpectedResponse expects a request for login and returns a 200 as an unexpected response.
func WithAuthPasswordLoginUnexpectedResponse(username, password string, deviceID uuid.UUID) ServerOption {
	return func(s *Server) {
		expectAuthPasswordLogin(s, username, password, deviceID).
			ReturnCode(http.StatusOK)
	}
}

// WithAuthPasswordLoginFailure expects a request for login and returns a 500.
func WithAuthPasswordLoginFailure(username, password string, deviceID uuid.UUID) ServerOption {
	return func(s *Server) {
		expectAuthPasswordLogin(s, username, password, deviceID).
			ReturnCode(http.StatusInternalServerError)
	}
}

// WithAuthPasswordLoginSuccess expects a request for login and returns a 403 for the MFA Challenge.
func WithAuthPasswordLoginSuccess(username, password string, deviceID uuid.UUID) ServerOption {
	return func(s *Server) {
		expectAuthPasswordLogin(s, username, password, deviceID).
			ReturnCode(http.StatusForbidden).
			Handler(func(r *http.Request) ([]byte, error) {
				mfaToken := uuid.New()

				s.WithMFAToken(mfaToken)

				response := api.RequiredMFATokenError{
					UserMessage: api.UserMessage{
						Title:  "A second authentication factor is required.",
						Detail: "Please provide your second form of authentication.",
					},
					MfaToken:         mfaToken.String(),
					ErrorDescription: "MFA token is required",
					Detail:           "MFA token is required",
					HostURL:          s.URL(),
					Type:             "mfa_required",
					Error:            "mfa_required",
					Title:            "A second authentication factor is required.",
					Message:          "Please provide your second form of authentication.",
					UserID:           s.userID.String(),
					Status:           http.StatusForbidden,
				}

				return json.Marshal(response)
			})
	}
}

// WithAuthMFAChallengeFailureInvalidToken expects a request for MFA Challenge and returns an Invalid Token error (401).
func WithAuthMFAChallengeFailureInvalidToken() ServerOption {
	return func(s *Server) {
		expectMFAChallenge(s).
			ReturnCode(http.StatusUnauthorized).
			ReturnJSON(api.InvalidTokenError{
				Status: http.StatusUnauthorized,
				Detail: "Invalid token",
				Type:   "error",
				UserMessage: api.UserMessage{
					Title:  "Login attempt expired",
					Detail: "That took too long, please try again.",
				},
				Error:            "invalid_token",
				ErrorDescription: "Invalid token",
			})
	}
}

// WithAuthMFAChallengeFailure expects a request for MFA Challenge and returns a 500.
func WithAuthMFAChallengeFailure() ServerOption {
	return func(s *Server) {
		expectMFAChallenge(s).ReturnCode(http.StatusInternalServerError)
	}
}

// WithAuthMFAChallengeSuccess expects a request for MFA Challenge and returns a success.
func WithAuthMFAChallengeSuccess() ServerOption {
	return func(s *Server) {
		expectMFAChallenge(s).
			ReturnCode(http.StatusCreated).
			ReturnJSON(api.PostAPIMfaChallengeResponseValueCreated{
				ChallengeType: "oob",
			})
	}
}

// WithAuthConfirmLoginFailureInvalidToken expects a request for Login Confirm and returns an Invalid Token error (401).
func WithAuthConfirmLoginFailureInvalidToken(times int) ServerOption {
	return func(s *Server) {
		expectConfirmLogin(s).
			ReturnCode(http.StatusUnauthorized).
			ReturnJSON(api.InvalidTokenError{
				Status: http.StatusUnauthorized,
				Detail: "Invalid token",
				Type:   "error",
				UserMessage: api.UserMessage{
					Title:  "Login attempt expired",
					Detail: "That took too long, please try again.",
				},
				Error:            "invalid_token",
				ErrorDescription: "Invalid token",
			}).
			Times(times)
	}
}

// WithAuthConfirmLoginFailure expects a request for Login Confirm and returns a 500.
func WithAuthConfirmLoginFailure() ServerOption {
	return func(s *Server) {
		expectConfirmLogin(s).
			ReturnCode(http.StatusInternalServerError)
	}
}

// WithAuthConfirmLoginSuccess expects a request for Login Confirm and returns a success.
func WithAuthConfirmLoginSuccess() ServerOption {
	return func(s *Server) {
		expectConfirmLogin(s).
			ReturnCode(http.StatusOK).
			Handler(func(r *http.Request) ([]byte, error) {
				accessToken := uuid.New()
				refreshToken := uuid.New()

				s.WithAccessToken(accessToken).
					WithRefreshToken(refreshToken)

				response := api.TokenResponse{
					AccessToken:  accessToken.String(),
					TokenType:    "bearer",
					RefreshToken: refreshToken.String(),
					ExpiresIn:    889,
					HostURL:      s.URL(),
				}

				return json.Marshal(response)
			})
	}
}

// WithAuthRefreshTokenFailureInvalidToken expects a request for Token Refresh and returns an Invalid Token error (401).
func WithAuthRefreshTokenFailureInvalidToken() ServerOption {
	return func(s *Server) {
		expectRefreshToken(s).
			ReturnCode(http.StatusUnauthorized).
			ReturnJSON(api.InvalidTokenError{
				Status: http.StatusUnauthorized,
				Detail: "Invalid token",
				Type:   "error",
				UserMessage: api.UserMessage{
					Title:  "Login attempt expired",
					Detail: "That took too long, please try again.",
				},
				Error:            "invalid_token",
				ErrorDescription: "Invalid token",
			})
	}
}

// WithAuthRefreshTokenFailure expects a request for Token Refresh and returns a 500.
func WithAuthRefreshTokenFailure() ServerOption {
	return func(s *Server) {
		expectRefreshToken(s).
			ReturnCode(http.StatusInternalServerError)
	}
}

// WithAuthRefreshTokenSuccess expects a request for Token Refresh and returns a success.
func WithAuthRefreshTokenSuccess() ServerOption {
	return func(s *Server) {
		expectRefreshToken(s).
			ReturnCode(http.StatusOK).
			Handler(func(r *http.Request) ([]byte, error) {
				accessToken := uuid.New()
				refreshToken := uuid.New()

				s.WithAccessToken(accessToken).
					WithRefreshToken(refreshToken)

				response := api.TokenResponse{
					AccessToken:  accessToken.String(),
					TokenType:    "bearer",
					RefreshToken: refreshToken.String(),
					ExpiresIn:    889,
					HostURL:      s.URL(),
				}

				return json.Marshal(response)
			})
	}
}

// WithAuthSuccess expects a success login workflow.
func WithAuthSuccess(username, password string, deviceID uuid.UUID) ServerOption {
	return func(s *Server) {
		WithAuthPasswordLoginSuccess(username, password, deviceID)(s)
		WithAuthMFAChallengeSuccess()(s)
		WithAuthConfirmLoginSuccess()(s)
	}
}
