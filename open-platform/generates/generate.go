package generates

import (
	"context"
	"net/http"
	"time"

	"github.com/767829413/advanced-go/open-platform/models"
)

type (
	// GenerateBasic provide the basis of the generated token data
	GenerateBasic struct {
		Client    models.ClientInfo
		UserID    string
		CreateAt  time.Time
		TokenInfo models.TokenInfo
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
)
