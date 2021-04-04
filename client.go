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
	credentials []CredentialsProvider
	tokens      []auth.TokenProvider
	transport   http.RoundTripper

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

	credentials := chainCredentialsProviders(
		CredentialsFromEnv(),
		Credentials(c.config.username, c.config.password),
	)

	for _, p := range c.config.credentials {
		credentials.chain(p)
	}

	token := chainTokenProviders(
		newAPITokenProvider(c.config.baseURL, c.config.timeout, credentials, c.config.deviceID, c.clock).
			WithMFATimeout(c.config.mfaTimeout).
			WithMFAWait(c.config.mfaWait).
			WithTransport(c.config.transport),
	)

	for _, p := range c.config.tokens {
		token.chain(p)
	}

	c.token = token

	// Initiates API client.
	apiClient := api.NewClient()
	apiClient.BaseURL = c.config.baseURL
	apiClient.Timeout = c.config.timeout

	c.api = apiClient

	c.setTransport(c.config.transport)

	return c
}

func (c *Client) setTransport(transport http.RoundTripper) {
	c.api.SetTransport(TokenRoundTripper(c.token, transport))
}
