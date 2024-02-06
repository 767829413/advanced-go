package errors

import "errors"

// New returns an error that formats as the given text.
var New = errors.New

// known errors
var (
	ErrInvalidRedirectURI   = errors.New("invalid redirect uri")
	ErrInvalidAuthorizeCode = errors.New("invalid authorize code")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrExpiredAccessToken   = errors.New("expired access token")
	ErrExpiredRefreshToken  = errors.New("expired refresh token")
	ErrMissingCodeVerifier  = errors.New("missing code verifier")
	ErrMissingCodeChallenge = errors.New("missing code challenge")
	ErrInvalidCodeChallenge = errors.New("invalid code challenge")
)

// https://tools.ietf.org/html/rfc6749#section-5.2
var (
	ErrInvalidRequest                 = errors.New("invalid_request")
	ErrUnauthorizedClient             = errors.New("unauthorized_client")
	ErrAccessDenied                   = errors.New("access_denied")
	ErrUnsupportedResponseType        = errors.New("unsupported_response_type")
	ErrInvalidScope                   = errors.New("invalid_scope")
	ErrServerError                    = errors.New("server_error")
	ErrTemporarilyUnavailable         = errors.New("temporarily_unavailable")
	ErrInvalidClient                  = errors.New("invalid_client")
	ErrInvalidGrant                   = errors.New("invalid_grant")
	ErrUnsupportedGrantType           = errors.New("unsupported_grant_type")
	ErrCodeChallengeRquired           = errors.New("invalid_request")
	ErrUnsupportedCodeChallengeMethod = errors.New("invalid_request")
	ErrInvalidCodeChallengeLen        = errors.New("invalid_request")
)
