package n26api

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/nhatthm/n26api/internal/api"
	"github.com/nhatthm/n26api/pkg/auth"
)

const (
	// BaseURL is N26 API Base URL.
	BaseURL = "https://api.tech26.de"

	envDeviceID = "N26_DEVICE"
)

var emptyUUID uuid.UUID

// Option configures Client.
type Option func(c *Client)

// Client provides all N26 APIs.
type Client struct {
	api   *api.Client
	token auth.TokenProvider
	clock Clock

	config *config
}

// config is configuration of Client.
type config struct {
	credentials  []CredentialsProvider
	tokens       []auth.TokenProvider
	tokenStorage auth.TokenStorage
	transport    http.RoundTripper

	baseURL  string
	timeout  time.Duration
	username string
	password string
	deviceID uuid.UUID

	mfaTimeout time.Duration
	mfaWait    time.Duration

	transactionsPageSize int64
}

// NewClient initiates a new transaction.Finder.
func NewClient(options ...Option) *Client {
	c := &Client{
		config: &config{
			transport: http.DefaultTransport,

			baseURL:  BaseURL,
			timeout:  time.Minute,
			deviceID: emptyUUID,

			mfaTimeout: time.Minute,
			mfaWait:    5 * time.Second,

			transactionsPageSize: transactionsPageSize,
		},

		clock: liveClock{},
	}

	for _, o := range options {
		o(c)
	}

	c.config.deviceID = deviceID(c.config.deviceID)
	c.token = initTokenProvider(c.config, c.clock)
	c.api = initAPIClient(c.config, c.token)

	return c
}

func initTokenProvider(cfg *config, c Clock) auth.TokenProvider {
	credentials := chainCredentialsProviders(
		CredentialsFromEnv(),
	)

	for _, p := range cfg.credentials {
		credentials.chain(p)
	}

	credentials.chain(Credentials(cfg.username, cfg.password))

	apiToken := newAPITokenProvider(credentials, cfg.deviceID).
		WithBaseURL(cfg.baseURL).
		WithTimeout(cfg.timeout).
		WithMFATimeout(cfg.mfaTimeout).
		WithMFAWait(cfg.mfaWait).
		WithTransport(cfg.transport).
		WithClock(c)

	if cfg.tokenStorage != nil {
		apiToken.WithStorage(cfg.tokenStorage)
	}

	token := chainTokenProviders(apiToken)

	for _, p := range cfg.tokens {
		token.chain(p)
	}

	return token
}

func initAPIClient(cfg *config, p auth.TokenProvider) *api.Client {
	c := api.NewClient()
	c.BaseURL = cfg.baseURL
	c.Timeout = cfg.timeout

	c.SetTransport(TokenRoundTripper(p, cfg.transport))

	return c
}
