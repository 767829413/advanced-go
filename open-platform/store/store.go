package store

import (
	"context"

	"github.com/767829413/advanced-go/open-platform/generates"
)

type (
	// ClientStore the client information storage interface
	ClientStore interface {
		// according to the ID for the client information
		GetByID(ctx context.Context, id string) (generates.ClientInfo, error)
	}

	// TokenStore the token information storage interface
	TokenStore interface {
		// create and store the new token information
		Create(ctx context.Context, info generates.TokenInfo) error

		// delete the authorization code
		RemoveByCode(ctx context.Context, code string) error

		// use the access token to delete the token information
		RemoveByAccess(ctx context.Context, access string) error

		// use the refresh token to delete the token information
		RemoveByRefresh(ctx context.Context, refresh string) error

		// use the authorization code for token information data
		GetByCode(ctx context.Context, code string) (generates.TokenInfo, error)

		// use the access token for token information data
		GetByAccess(ctx context.Context, access string) (generates.TokenInfo, error)

		// use the refresh token for token information data
		GetByRefresh(ctx context.Context, refresh string) (generates.TokenInfo, error)
	}
)