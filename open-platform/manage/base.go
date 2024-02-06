package manage

import (
	"context"
	"github.com/767829413/advanced-go/open-platform/config"
	"github.com/767829413/advanced-go/open-platform/generates"
	"net/http"
	"time"
)

// TokenGenerateRequest provide to generate the token request parameters
type TokenGenerateRequest struct {
	ClientID            string
	ClientSecret        string
	UserID              string
	RedirectURI         string
	Scope               string
	Code                string
	CodeChallenge       string
	CodeChallengeMethod config.CodeChallengeMethod
	Refresh             string
	CodeVerifier        string
	AccessTokenExp      time.Duration
	Request             *http.Request
}

// Manager authorization management interface
type Manager interface {

	// generate the authorization token(code)
	GenerateAuthToken(
		ctx context.Context,
		rt config.ResponseType,
		tgr *TokenGenerateRequest,
		cli generates.ClientInfo,
	) (authToken generates.TokenInfo, err error)

	// generate the access token
	GenerateAccessToken(
		ctx context.Context,
		gt config.GrantType,
		tgr *TokenGenerateRequest,
		cli generates.ClientInfo,
	) (accessToken generates.TokenInfo, err error)

	// refreshing an access token
	RefreshAccessToken(
		ctx context.Context,
		tgr *TokenGenerateRequest,
		cli generates.ClientInfo,
	) (accessToken generates.TokenInfo, err error)

	// use the access token to delete the token information
	RemoveAccessToken(ctx context.Context, access string) (err error)

	// use the refresh token to delete the token information
	RemoveRefreshToken(ctx context.Context, refresh string) (err error)

	// according to the access token for corresponding token information
	LoadAccessToken(ctx context.Context, access string) (ti generates.TokenInfo, err error)

	// according to the refresh token for corresponding token information
	LoadRefreshToken(ctx context.Context, refresh string) (ti generates.TokenInfo, err error)
}
