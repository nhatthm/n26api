package n26api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nhatthm/go-clock"

	"github.com/nhatthm/n26api/internal/api"
	"github.com/nhatthm/n26api/pkg/auth"
)

const (
	// BaseURL is N26 API Base URL.
	BaseURL = api.DefaultBaseURL
	// DefaultPageSize is the default page size while requesting to N26.
	DefaultPageSize int64 = 50

	envDeviceID = "N26_DEVICE"
)

var emptyUUID uuid.UUID

// Option configures Client.
type Option func(c *Client)

// Client provides all N26 APIs.
type Client struct {
	api   *api.Client
	token *chainTokenProvider
	clock clock.Clock

	config *config
}

// config is configuration of Client.
type config struct {
	credentials  *chainCredentialsProvider
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

// DeviceID returns device ID.
func (c *Client) DeviceID() uuid.UUID {
	return c.config.deviceID
}

// NewClient initiates a new transaction.Finder.
func NewClient(options ...Option) *Client {
	c := &Client{
		config: &config{
			credentials: chainCredentialsProviders(CredentialsFromEnv()),
			transport:   http.DefaultTransport,

			baseURL:  BaseURL,
			timeout:  time.Minute,
			deviceID: emptyUUID,

			mfaTimeout: time.Minute,
			mfaWait:    5 * time.Second,

			transactionsPageSize: DefaultPageSize,
		},

		token: newChainTokenProvider(),
		clock: clock.New(),
	}

	for _, o := range options {
		o(c)
	}

	c.config.deviceID = deviceID(c.config.deviceID)
	c.token.append(initAPITokenProvider(c.config, c.clock))
	c.api = initAPIClient(c.config, c.token)

	return c
}

func initAPITokenProvider(cfg *config, c clock.Clock) auth.TokenProvider {
	cfg.credentials.prepend(Credentials(cfg.username, cfg.password))

	apiToken := newAPITokenProvider(cfg.credentials, cfg.deviceID).
		WithBaseURL(cfg.baseURL).
		WithTimeout(cfg.timeout).
		WithMFATimeout(cfg.mfaTimeout).
		WithMFAWait(cfg.mfaWait).
		WithTransport(cfg.transport).
		WithClock(c)

	if cfg.tokenStorage != nil {
		apiToken.WithStorage(cfg.tokenStorage)
	}

	return apiToken
}

func initAPIClient(cfg *config, p auth.TokenProvider) *api.Client {
	c := api.NewClient()
	c.BaseURL = cfg.baseURL
	c.Timeout = cfg.timeout

	c.SetTransport(TokenRoundTripper(p, cfg.transport))

	return c
}
