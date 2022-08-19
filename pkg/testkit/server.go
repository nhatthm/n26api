package testkit

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/nhatthm/httpmock"
	"github.com/nhatthm/httpmock/matcher"
	"github.com/nhatthm/httpmock/planner"
	"github.com/nhatthm/httpmock/request"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/nhatthm/n26api/pkg/util"
)

// Request is an alias of httpmock.Request.
type Request = request.Request

// Server is a wrapped httpmock.Server to provide more functionalities for testing N26 APIs.
type Server struct {
	*httpmock.Server

	authUsername string
	authPassword string
	userID       uuid.UUID
	deviceID     uuid.UUID
	mfaToken     uuid.UUID
	accessToken  auth.Token
	refreshToken auth.Token

	mu sync.Mutex
}

// WithPlanner sets planner.
func (s *Server) WithPlanner(p planner.Planner) *Server {
	s.Server.WithPlanner(p)

	return s
}

// WithAuthAuthorization sets Authorization credentials for asserting the /oauth/token request.
func (s *Server) WithAuthAuthorization(username, password string) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.authUsername = username
	s.authPassword = password

	return s
}

// WithDeviceID sets the deviceID.
func (s *Server) WithDeviceID(deviceID uuid.UUID) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.deviceID = deviceID

	return s
}

// WithMFAToken sets the mfaToken.
func (s *Server) WithMFAToken(mfaToken uuid.UUID) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.mfaToken = mfaToken

	return s
}

// WithAccessToken sets the accessToken.
func (s *Server) WithAccessToken(token uuid.UUID) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.accessToken = auth.Token(token.String())

	return s
}

// WithRefreshToken sets the refreshToken.
func (s *Server) WithRefreshToken(token uuid.UUID) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.refreshToken = auth.Token(token.String())

	return s
}

// DeviceID returns the deviceID.
func (s *Server) DeviceID() uuid.UUID {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.deviceID
}

// MFAToken returns the deviceID.
func (s *Server) MFAToken() uuid.UUID {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.mfaToken
}

// BasicAuthorization returns the credentials for asserting the /oauth/token Authorization.
func (s *Server) BasicAuthorization() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.authUsername == "" && s.authPassword == "" {
		return "Basic"
	}

	return fmt.Sprintf("Basic %s", util.Base64Credentials(s.authUsername, s.authPassword))
}

// UserID returns the userID.
func (s *Server) UserID() uuid.UUID {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.userID
}

// AccessToken returns the accessToken.
func (s *Server) AccessToken() auth.Token {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.accessToken
}

// RefreshToken returns the refreshToken.
func (s *Server) RefreshToken() auth.Token {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.refreshToken
}

// ExpectWithBasicAuth expects a request with Basic Authorization.
func (s *Server) ExpectWithBasicAuth(method string, requestURI interface{}) *Request {
	return s.Server.Expect(method, requestURI).
		WithHeader("Authorization", func() matcher.Matcher {
			return httpmock.Exact(s.BasicAuthorization())
		})
}

// Expect expects a request with Bearer Authorization.
//
//	Server.Expect(http.MethodGet, "/path").
func (s *Server) Expect(method string, requestURI interface{}) *Request {
	return s.Server.Expect(method, requestURI).
		WithHeader("Authorization", func() matcher.Matcher {
			return httpmock.Exactf("Bearer %s", s.accessToken)
		})
}

// ExpectGet expects a request with Bearer Authorization.
//
//	Server.ExpectGet("/path")
func (s *Server) ExpectGet(requestURI interface{}) *Request {
	return s.Expect(http.MethodGet, requestURI)
}

// ExpectHead expects a request with Bearer Authorization.
//
//	Server.ExpectHead("/path")
func (s *Server) ExpectHead(requestURI interface{}) *Request {
	return s.Expect(http.MethodHead, requestURI)
}

// ExpectPost expects a request with Bearer Authorization.
//
//	Server.ExpectPost("/path")
func (s *Server) ExpectPost(requestURI interface{}) *Request {
	return s.Expect(http.MethodPost, requestURI)
}

// ExpectPut expects a request with Bearer Authorization.
//
//	Server.ExpectPut("/path")
func (s *Server) ExpectPut(requestURI interface{}) *Request {
	return s.Expect(http.MethodPut, requestURI)
}

// ExpectPatch expects a request with Bearer Authorization.
//
//	Server.ExpectPatch("/path")
func (s *Server) ExpectPatch(requestURI interface{}) *Request {
	return s.Expect(http.MethodPatch, requestURI)
}

// ExpectDelete expects a request with Bearer Authorization.
//
//	Server.ExpectDelete("/path")
func (s *Server) ExpectDelete(requestURI interface{}) *Request {
	return s.Expect(http.MethodDelete, requestURI)
}

// NewServer creates a new Server.
func NewServer(t TestingT) *Server {
	s := &Server{
		Server: httpmock.NewServer().WithTest(t),
		userID: uuid.New(),
	}

	s.WithAuthAuthorization(auth.BasicAuthUsername, auth.BasicAuthPassword).
		WithDefaultResponseHeaders(httpmock.Header{
			"Content-Type": "application/json",
		})

	return s
}
