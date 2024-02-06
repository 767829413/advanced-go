package config

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"
)

// Config authorization configuration parameters
type Config struct {
	// access token expiration time, 0 means it doesn't expire
	AccessTokenExp time.Duration
	// refresh token expiration time, 0 means it doesn't expire
	RefreshTokenExp time.Duration
	// whether to generate the refreshing token
	IsGenerateRefresh bool
}

// RefreshingConfig refreshing token config
type RefreshingConfig struct {
	// access token expiration time, 0 means it doesn't expire
	AccessTokenExp time.Duration
	// refresh token expiration time, 0 means it doesn't expire
	RefreshTokenExp time.Duration
	// whether to generate the refreshing token
	IsGenerateRefresh bool
	// whether to reset the refreshing create time
	IsResetRefreshTime bool
	// whether to remove access token
	IsRemoveAccess bool
	// whether to remove refreshing token
	IsRemoveRefreshing bool
}

// default configs
var (
	DefaultCodeExp               = time.Minute * 10
	DefaultAuthorizeCodeTokenCfg = &Config{
		AccessTokenExp:    time.Hour * 2,
		RefreshTokenExp:   time.Hour * 24 * 3,
		IsGenerateRefresh: true,
	}
	DefaultImplicitTokenCfg = &Config{AccessTokenExp: time.Hour * 1}
	DefaultPasswordTokenCfg = &Config{
		AccessTokenExp:    time.Hour * 2,
		RefreshTokenExp:   time.Hour * 24 * 7,
		IsGenerateRefresh: true,
	}
	DefaultClientTokenCfg  = &Config{AccessTokenExp: time.Hour * 2}
	DefaultRefreshTokenCfg = &RefreshingConfig{
		IsGenerateRefresh:  true,
		IsRemoveAccess:     true,
		IsRemoveRefreshing: true,
	}
)

// GrantType authorization model
type GrantType string

// define authorization model
const (
	AuthorizationCode   GrantType = "authorization_code"
	PasswordCredentials GrantType = "password"
	ClientCredentials   GrantType = "client_credentials"
	Refreshing          GrantType = "refresh_token"
	Implicit            GrantType = "__implicit"
)

// ResponseType the type of authorization request
type ResponseType string

// define the type of authorization request
const (
	Code  ResponseType = "code"
	Token ResponseType = "token"
)

func (gt GrantType) String() string {
	if gt == AuthorizationCode ||
		gt == PasswordCredentials ||
		gt == ClientCredentials ||
		gt == Refreshing {
		return string(gt)
	}
	return ""
}

// CodeChallengeMethod PCKE method
type CodeChallengeMethod string

const (
	// CodeChallengePlain PCKE Method
	CodeChallengePlain CodeChallengeMethod = "plain"
	// CodeChallengeS256 PCKE Method
	CodeChallengeS256 CodeChallengeMethod = "S256"
)

func (ccm CodeChallengeMethod) String() string {
	if ccm == CodeChallengePlain ||
		ccm == CodeChallengeS256 {
		return string(ccm)
	}
	return ""
}

// Validate code challenge
func (ccm CodeChallengeMethod) Validate(cc, ver string) bool {
	switch ccm {
	case CodeChallengePlain:
		return cc == ver
	case CodeChallengeS256:
		s256 := sha256.Sum256([]byte(ver))
		// trim padding
		a := strings.TrimRight(base64.URLEncoding.EncodeToString(s256[:]), "=")
		b := strings.TrimRight(cc, "=")
		return a == b
	default:
		return false
	}
}
