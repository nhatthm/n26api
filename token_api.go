package n26api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bool64/ctxd"
	"github.com/google/uuid"
	"go.nhat.io/clock"

	"github.com/nhatthm/n26api/internal/api"
	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/nhatthm/n26api/pkg/util"
)

var (
	// ErrUsernameIsEmpty indicates that the username is empty.
	ErrUsernameIsEmpty = errors.New("missing username")
	// ErrPasswordIsEmpty indicates that the username is empty.
	ErrPasswordIsEmpty = errors.New("missing password")
)

var _ auth.TokenProvider = (*apiTokenProvider)(nil)

var emptyToken = auth.OAuthToken{}

type apiTokenProvider struct {
	api         *api.Client
	credentials CredentialsProvider
	storage     auth.TokenStorage
	clock       clock.Clock

	deviceID uuid.UUID

	mfaTimeout time.Duration
	mfaWait    time.Duration
	refreshTTL time.Duration

	mu sync.Mutex
}

func (p *apiTokenProvider) getToken(ctx context.Context, key string) (auth.OAuthToken, error) {
	return p.storage.Get(ctx, key)
}

func (p *apiTokenProvider) setToken(ctx context.Context, key string, res api.TokenResponse, timestamp time.Time) (auth.OAuthToken, error) {
	expiryDuration := time.Duration(res.ExpiresIn * int64(time.Second))

	token := auth.OAuthToken{
		AccessToken:      auth.Token(res.AccessToken),
		RefreshToken:     auth.Token(res.RefreshToken),
		ExpiresAt:        timestamp.Add(expiryDuration),
		RefreshExpiresAt: timestamp.Add(p.refreshTTL),
	}

	if err := p.storage.Set(ctx, key, token); err != nil {
		return auth.OAuthToken{}, err
	}

	return token, nil
}

func (p *apiTokenProvider) login(ctx context.Context) (string, error) {
	password := p.credentials.Password()
	if password == "" {
		return "", ctxd.WrapError(ctx, ErrPasswordIsEmpty, "could not get token")
	}

	res, err := p.api.PostOauthToken(ctx, api.PostOauthTokenRequest{
		DeviceToken: p.deviceID.String(),
		GrantType:   "password",
		Username:    util.StringPtr(p.credentials.Username()),
		Password:    util.StringPtr(p.credentials.Password()),
	})
	if err != nil {
		return "", ctxd.WrapError(ctx, err, "unexpected response")
	}

	statusCode := res.StatusCode

	switch statusCode {
	case http.StatusBadRequest:
		return "", ctxd.NewError(ctx, "wrong credentials", "response", res)

	case http.StatusForbidden:
		return res.ValueForbidden.MfaToken, nil

	case http.StatusTooManyRequests:
		return "", ctxd.NewError(ctx, "too many login attempts", "response", res)
	}

	return "", err
}

func (p *apiTokenProvider) challenge(ctx context.Context, token string) error {
	res, err := p.api.PostAPIMfaChallenge(ctx, api.PostAPIMfaChallengeRequest{
		DeviceToken: p.deviceID.String(),
		Body: &api.MFAChallengeRequest{
			ChallengeType: "oob",
			MfaToken:      token,
		},
	})
	if err != nil {
		return ctxd.WrapError(ctx, err, "failed to challenge mfa")
	}

	if res.ValueCreated == nil {
		return ctxd.NewError(ctx, "could not challenge mfa", "response", res)
	}

	return nil
}

func (p *apiTokenProvider) confirmLogin(ctx context.Context, token string) (*api.TokenResponse, error) {
	res, err := p.api.PostOauthToken(ctx, api.PostOauthTokenRequest{
		DeviceToken: p.deviceID.String(),
		GrantType:   "mfa_oob",
		MfaToken:    util.StringPtr(token),
	})
	if err != nil {
		return nil, ctxd.WrapError(ctx, err, "failed to confirm login")
	}

	if res.ValueOK == nil {
		return nil, ctxd.NewError(ctx, "could not get access token", "response", res)
	}

	return res.ValueOK, nil
}

func (p *apiTokenProvider) get(ctx context.Context, key string, timestamp time.Time) (auth.Token, error) {
	mfaToken, err := p.login(ctx)
	if err != nil {
		return "", err
	}

	if err := p.challenge(ctx, mfaToken); err != nil {
		return "", err
	}

	timeout, cancel := context.WithTimeout(ctx, p.mfaTimeout)
	defer cancel()

	ticker := time.NewTicker(p.mfaWait)

	for {
		select {
		case <-ticker.C:
			res, _ := p.confirmLogin(timeout, mfaToken) // nolint:errcheck

			if res != nil {
				token, err := p.setToken(ctx, key, *res, timestamp)
				if err != nil {
					return "", ctxd.WrapError(ctx, err, "could not persist token to storage")
				}

				return token.AccessToken, nil
			}

		case <-timeout.Done():
			return "", ctxd.NewError(ctx, "could not confirm login", "reason", "timeout")
		}
	}
}

func (p *apiTokenProvider) refresh(ctx context.Context, key string, refreshToken auth.Token, timestamp time.Time) (auth.Token, error) {
	res, err := p.api.PostOauthToken(ctx, api.PostOauthTokenRequest{
		DeviceToken:  p.deviceID.String(),
		GrantType:    "refresh_token",
		RefreshToken: util.StringPtr(string(refreshToken)),
	})
	if err != nil {
		return "", ctxd.WrapError(ctx, err, "failed to refresh token")
	}

	if res.ValueOK != nil {
		token, err := p.setToken(ctx, key, *res.ValueOK, timestamp)
		if err != nil {
			return "", ctxd.WrapError(ctx, err, "could not persist token to storage")
		}

		return token.AccessToken, nil
	}

	return p.get(ctx, key, timestamp)
}

func (p *apiTokenProvider) WithBaseURL(baseURL string) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.api.BaseURL = baseURL

	return p
}

func (p *apiTokenProvider) WithTimeout(timeout time.Duration) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.api.Timeout = timeout

	return p
}

// nolint:unparam
func (p *apiTokenProvider) WithStorage(storage auth.TokenStorage) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.storage = storage

	return p
}

func (p *apiTokenProvider) WithTransport(transport http.RoundTripper) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.api.SetTransport(BasicAuthRoundTripper(auth.BasicAuthUsername, auth.BasicAuthPassword, transport))

	return p
}

func (p *apiTokenProvider) WithMFATimeout(timeout time.Duration) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mfaTimeout = timeout

	return p
}

func (p *apiTokenProvider) WithMFAWait(waitTime time.Duration) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mfaWait = waitTime

	return p
}

func (p *apiTokenProvider) WithRefreshTTL(ttl time.Duration) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.refreshTTL = ttl

	return p
}

func (p *apiTokenProvider) WithClock(clock clock.Clock) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clock = clock

	return p
}

func (p *apiTokenProvider) Token(ctx context.Context) (auth.Token, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	username := p.credentials.Username()
	if username == "" {
		return "", ctxd.WrapError(ctx, ErrUsernameIsEmpty, "could not get token")
	}

	key := fmt.Sprintf("%s:%s", username, p.deviceID.String())
	now := p.clock.Now()

	token, err := p.getToken(ctx, key)
	if err != nil {
		return "", ctxd.WrapError(ctx, err, "could not get token from storage")
	}

	if token == emptyToken {
		return p.get(ctx, key, now)
	}

	if !token.IsExpired(now) {
		return token.AccessToken, nil
	}

	if token.IsRefreshable(now) {
		return p.refresh(ctx, key, token.RefreshToken, now)
	}

	return p.get(ctx, key, now)
}

func newAPITokenProvider(
	credentials CredentialsProvider,
	deviceID uuid.UUID,
) *apiTokenProvider {
	c := api.NewClient()
	c.BaseURL = BaseURL
	c.Timeout = time.Minute
	c.SetTransport(BasicAuthRoundTripper(
		auth.BasicAuthUsername, auth.BasicAuthPassword,
		http.DefaultTransport,
	))

	return &apiTokenProvider{
		api:         c,
		credentials: credentials,
		storage:     NewInMemoryTokenStorage(),
		clock:       clock.New(),

		deviceID: deviceID,

		mfaTimeout: time.Minute,
		mfaWait:    5 * time.Second,
		refreshTTL: time.Hour,
	}
}
