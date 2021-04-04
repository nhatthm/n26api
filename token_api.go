package n26api

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/bool64/ctxd"
	"github.com/google/uuid"

	"github.com/nhatthm/n26api/internal/api"
	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/nhatthm/n26api/pkg/util"
)

var _ auth.TokenProvider = (*apiTokenProvider)(nil)

var emptyToken = apiToken{}

type apiTokenProvider struct {
	api         *api.Client
	credentials CredentialsProvider
	clock       Clock

	deviceID uuid.UUID
	token    apiToken

	mfaTimeout time.Duration
	mfaWait    time.Duration
	refreshTTL time.Duration

	mu sync.Mutex
}

type apiToken struct {
	accessToken      auth.Token
	refreshToken     auth.Token
	expiresAt        time.Time
	refreshExpiresAt time.Time
}

func (p *apiTokenProvider) getToken() apiToken {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.token
}

func (p *apiTokenProvider) setToken(res api.TokenResponse, timestamp time.Time) apiToken {
	p.mu.Lock()
	defer p.mu.Unlock()

	expiryDuration := time.Duration(res.ExpiresIn * int64(time.Second))

	p.token = apiToken{
		accessToken:      auth.Token(res.AccessToken),
		refreshToken:     auth.Token(res.RefreshToken),
		expiresAt:        timestamp.Add(expiryDuration),
		refreshExpiresAt: timestamp.Add(p.refreshTTL),
	}

	return p.token
}

func (p *apiTokenProvider) login(ctx context.Context) (string, error) {
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

func (p *apiTokenProvider) get(ctx context.Context, timestamp time.Time) (auth.Token, error) {
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
				return p.setToken(*res, timestamp).accessToken, nil
			}

		case <-timeout.Done():
			return "", ctxd.NewError(ctx, "could not confirm login", "reason", "timeout")
		}
	}
}

func (p *apiTokenProvider) refresh(ctx context.Context, timestamp time.Time) (auth.Token, error) {
	res, err := p.api.PostOauthToken(ctx, api.PostOauthTokenRequest{
		DeviceToken:  p.deviceID.String(),
		GrantType:    "refresh_token",
		RefreshToken: util.StringPtr(string(p.getToken().refreshToken)),
	})
	if err != nil {
		return "", ctxd.WrapError(ctx, err, "failed to refresh token")
	}

	if res.ValueOK != nil {
		return p.setToken(*res.ValueOK, timestamp).accessToken, nil
	}

	return p.get(ctx, timestamp)
}

func (p *apiTokenProvider) WithTransport(transport http.RoundTripper) *apiTokenProvider {
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

func (p *apiTokenProvider) Token(ctx context.Context) (auth.Token, error) {
	now := p.clock.Now()
	token := p.getToken()

	if token == emptyToken {
		return p.get(ctx, now)
	}

	if !token.isExpired(now) {
		return token.accessToken, nil
	}

	if token.isRefreshable(now) {
		return p.refresh(ctx, now)
	}

	return p.get(ctx, now)
}

func (t apiToken) isExpired(timestamp time.Time) bool {
	return t.expiresAt.Before(timestamp)
}

func (t apiToken) isRefreshable(timestamp time.Time) bool {
	return t.refreshExpiresAt.After(timestamp)
}

func newAPITokenProvider(
	baseURL string,
	timeout time.Duration,
	credentials CredentialsProvider,
	deviceID uuid.UUID,
	clock Clock,
) *apiTokenProvider {
	c := api.NewClient()
	c.BaseURL = baseURL
	c.Timeout = timeout
	c.SetTransport(BasicAuthRoundTripper(
		auth.BasicAuthUsername, auth.BasicAuthPassword,
		http.DefaultTransport,
	))

	return &apiTokenProvider{
		api:         c,
		credentials: credentials,
		clock:       clock,

		deviceID: deviceID,
		token:    emptyToken,

		mfaTimeout: time.Minute,
		mfaWait:    5 * time.Second,
		refreshTTL: time.Hour,
	}
}
