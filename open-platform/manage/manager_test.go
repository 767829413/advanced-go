package manage_test

import (
	"context"
	"testing"
	"time"

	"github.com/767829413/advanced-go/open-platform/config"
	"github.com/767829413/advanced-go/open-platform/manage"
	"github.com/767829413/advanced-go/open-platform/models"
	"github.com/767829413/advanced-go/open-platform/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	Convey("Manager test", t, func() {
		manager := manage.NewDefaultManager()
		ctx := context.Background()

		manager.MustTokenStorage(store.NewMemoryTokenStore())

		clientStore := store.NewClientStore()
		_ = clientStore.Set("1", &models.Client{
			ID:     "1",
			Secret: "11",
			Domain: "http://localhost",
		})
		manager.MapClientStorage(clientStore)

		tgr := &manage.TokenGenerateRequest{
			ClientID:    "1",
			UserID:      "123456",
			RedirectURI: "http://localhost/oauth2",
			Scope:       "all",
		}

		Convey("GetClient test", func() {
			cli, err := manager.GetClient(ctx, "1")
			So(err, ShouldBeNil)
			So(cli.GetSecret(), ShouldEqual, "11")
		})

		Convey("Token test", func() {
			testManager(tgr, manager)
		})

		Convey("zero expiration access token test", func() {
			testZeroAccessExpirationManager(tgr, manager)
			testCannotRequestZeroExpirationAccessTokens(tgr, manager)
		})

		Convey("zero expiration refresh token test", func() {
			testZeroRefreshExpirationManager(tgr, manager)
		})
	})
}

func testManager(tgr *manage.TokenGenerateRequest, manager manage.Manager) {
	ctx := context.Background()
	cti, err := manager.GenerateAuthToken(ctx, config.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &manage.TokenGenerateRequest{
		ClientID:     tgr.ClientID,
		ClientSecret: "11",
		RedirectURI:  tgr.RedirectURI,
		Code:         code,
	}
	ati, err := manager.GenerateAccessToken(ctx, config.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken, refreshToken := ati.GetAccess(), ati.GetRefresh()
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	ainfo, err := manager.LoadAccessToken(ctx, accessToken)
	So(err, ShouldBeNil)
	So(ainfo.GetClientID(), ShouldEqual, atParams.ClientID)

	arinfo, err := manager.LoadRefreshToken(ctx, accessToken)
	So(err, ShouldNotBeNil)
	So(arinfo, ShouldBeNil)

	rainfo, err := manager.LoadAccessToken(ctx, refreshToken)
	So(err, ShouldNotBeNil)
	So(rainfo, ShouldBeNil)

	rinfo, err := manager.LoadRefreshToken(ctx, refreshToken)
	So(err, ShouldBeNil)
	So(rinfo.GetClientID(), ShouldEqual, atParams.ClientID)

	refreshParams := &manage.TokenGenerateRequest{
		Refresh: refreshToken,
		Scope:   "owner",
	}
	rti, err := manager.RefreshAccessToken(ctx, refreshParams)
	So(err, ShouldBeNil)

	refreshAT := rti.GetAccess()
	So(refreshAT, ShouldNotBeEmpty)

	_, err = manager.LoadAccessToken(ctx, accessToken)
	So(err, ShouldNotBeNil)

	refreshAInfo, err := manager.LoadAccessToken(ctx, refreshAT)
	So(err, ShouldBeNil)
	So(refreshAInfo.GetScope(), ShouldEqual, "owner")

	err = manager.RemoveAccessToken(ctx, refreshAT)
	So(err, ShouldBeNil)

	_, err = manager.LoadAccessToken(ctx, refreshAT)
	So(err, ShouldNotBeNil)

	err = manager.RemoveRefreshToken(ctx, refreshToken)
	So(err, ShouldBeNil)

	_, err = manager.LoadRefreshToken(ctx, refreshToken)
	So(err, ShouldNotBeNil)
}

func testZeroAccessExpirationManager(tgr *manage.TokenGenerateRequest, manager manage.Manager) {
	ctx := context.Background()
	conf := config.Config{
		AccessTokenExp:    0, // Set explicitly as we're testing 0 (no) expiration
		IsGenerateRefresh: true,
	}
	m, ok := manager.(*manage.ManagerIns)
	So(ok, ShouldBeTrue)
	m.SetAuthorizeCodeTokenCfg(&conf)

	cti, err := manager.GenerateAuthToken(ctx, config.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &manage.TokenGenerateRequest{
		ClientID:     tgr.ClientID,
		ClientSecret: "11",
		RedirectURI:  tgr.RedirectURI,
		Code:         code,
	}
	ati, err := manager.GenerateAccessToken(ctx, config.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken, refreshToken := ati.GetAccess(), ati.GetRefresh()
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	tokenInfo, err := manager.LoadAccessToken(ctx, accessToken)
	So(err, ShouldBeNil)
	So(tokenInfo, ShouldNotBeNil)
	So(tokenInfo.GetAccess(), ShouldEqual, accessToken)
	So(tokenInfo.GetAccessExpiresIn(), ShouldEqual, 0)
}

func testCannotRequestZeroExpirationAccessTokens(
	tgr *manage.TokenGenerateRequest,
	manager manage.Manager,
) {
	ctx := context.Background()
	conf := config.Config{
		AccessTokenExp: time.Hour * 5,
	}
	m, ok := manager.(*manage.ManagerIns)
	So(ok, ShouldBeTrue)
	m.SetAuthorizeCodeTokenCfg(&conf)

	cti, err := manager.GenerateAuthToken(ctx, config.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &manage.TokenGenerateRequest{
		ClientID:       tgr.ClientID,
		ClientSecret:   "11",
		RedirectURI:    tgr.RedirectURI,
		AccessTokenExp: 0, // requesting token without expiration
		Code:           code,
	}
	ati, err := manager.GenerateAccessToken(ctx, config.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken := ati.GetAccess()
	So(accessToken, ShouldNotBeEmpty)
	So(ati.GetAccessExpiresIn(), ShouldEqual, time.Hour*5)
}

func testZeroRefreshExpirationManager(tgr *manage.TokenGenerateRequest, manager manage.Manager) {
	ctx := context.Background()
	conf := config.Config{
		RefreshTokenExp:   0, // Set explicitly as we're testing 0 (no) expiration
		IsGenerateRefresh: true,
	}
	m, ok := manager.(*manage.ManagerIns)
	So(ok, ShouldBeTrue)
	m.SetAuthorizeCodeTokenCfg(&conf)

	cti, err := manager.GenerateAuthToken(ctx, config.Code, tgr)
	So(err, ShouldBeNil)

	code := cti.GetCode()
	So(code, ShouldNotBeEmpty)

	atParams := &manage.TokenGenerateRequest{
		ClientID:       tgr.ClientID,
		ClientSecret:   "11",
		RedirectURI:    tgr.RedirectURI,
		AccessTokenExp: time.Hour,
		Code:           code,
	}
	ati, err := manager.GenerateAccessToken(ctx, config.AuthorizationCode, atParams)
	So(err, ShouldBeNil)

	accessToken, refreshToken := ati.GetAccess(), ati.GetRefresh()
	So(accessToken, ShouldNotBeEmpty)
	So(refreshToken, ShouldNotBeEmpty)

	tokenInfo, err := manager.LoadRefreshToken(ctx, refreshToken)
	So(err, ShouldBeNil)
	So(tokenInfo, ShouldNotBeNil)
	So(tokenInfo.GetRefresh(), ShouldEqual, refreshToken)
	So(tokenInfo.GetRefreshExpiresIn(), ShouldEqual, 0)

	// LoadAccessToken also checks refresh expiry
	tokenInfo, err = manager.LoadAccessToken(ctx, accessToken)
	So(err, ShouldBeNil)
	So(tokenInfo, ShouldNotBeNil)
	So(tokenInfo.GetRefresh(), ShouldEqual, refreshToken)
	So(tokenInfo.GetRefreshExpiresIn(), ShouldEqual, 0)
}
