package generates

import (
	"context"
	"net/http"
	"time"
)

type (
	// CodeChallengeMethod PCKE method
	CodeChallengeMethod string

	// GenerateBasic provide the basis of the generated token data
	GenerateBasic struct {
		Client    ClientInfo
		UserID    string
		CreateAt  time.Time
		TokenInfo TokenInfo
		Request   *http.Request
	}

	// AuthorizeGenerate generate the authorization code interface
	AuthorizeGenerate interface {
		Token(ctx context.Context, data *GenerateBasic) (code string, err error)
	}

	// AccessGenerate generate the access and refresh tokens interface
	AccessGenerate interface {
		Token(
			ctx context.Context,
			data *GenerateBasic,
			isGenRefresh bool,
		) (access, refresh string, err error)
	}

	// TokenInfo the token information model interface
	TokenInfo interface {
		New() TokenInfo

		GetClientID() string
		SetClientID(string)
		GetUserID() string
		SetUserID(string)
		GetRedirectURI() string
		SetRedirectURI(string)
		GetScope() string
		SetScope(string)

		GetCode() string
		SetCode(string)
		GetCodeCreateAt() time.Time
		SetCodeCreateAt(time.Time)
		GetCodeExpiresIn() time.Duration
		SetCodeExpiresIn(time.Duration)
		GetCodeChallenge() string
		SetCodeChallenge(string)
		GetCodeChallengeMethod() CodeChallengeMethod
		SetCodeChallengeMethod(CodeChallengeMethod)

		GetAccess() string
		SetAccess(string)
		GetAccessCreateAt() time.Time
		SetAccessCreateAt(time.Time)
		GetAccessExpiresIn() time.Duration
		SetAccessExpiresIn(time.Duration)

		GetRefresh() string
		SetRefresh(string)
		GetRefreshCreateAt() time.Time
		SetRefreshCreateAt(time.Time)
		GetRefreshExpiresIn() time.Duration
		SetRefreshExpiresIn(time.Duration)
	}

	// ClientInfo the client information model interface
	ClientInfo interface {
		GetID() string
		GetSecret() string
		GetDomain() string
		IsPublic() bool
		GetUserID() string
	}
)
