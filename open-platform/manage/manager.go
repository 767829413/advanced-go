package manage

import (
	"context"
	"time"

	"github.com/767829413/advanced-go/open-platform/config"
	"github.com/767829413/advanced-go/open-platform/errors"
	"github.com/767829413/advanced-go/open-platform/generates"
	"github.com/767829413/advanced-go/open-platform/models"
	"github.com/767829413/advanced-go/open-platform/store"
)

type ManagerIns struct {
	codeExp           time.Duration
	gtcfg             map[config.GrantType]*config.TokenExpConf
	rcfg              *config.RefreshingConfig
	authorizeGenerate generates.AuthorizeGenerate
	accessGenerate    generates.AccessGenerate
	tokenStore        store.TokenStore
	clientStore       store.ClientStore
	validateURI       ValidateURIHandler
}

// NewDefaultManager create to default authorization management instance
func NewDefaultManager() *ManagerIns {
	m := NewManager(nil)
	// default implementation
	m.MapAuthorizeGenerate(generates.NewAuthorizeGenerateIns())
	m.MapAccessGenerate(generates.NewAccessGenerate())

	return m
}

// NewManager create to authorization management instance
func NewManager(validateURI ValidateURIHandler) *ManagerIns {
	if validateURI == nil {
		validateURI = DefaultValidateURI
	}
	return &ManagerIns{
		gtcfg:       make(map[config.GrantType]*config.TokenExpConf),
		validateURI: validateURI,
	}
}

// get grant type config
func (m *ManagerIns) grantConfig(gt config.GrantType) *config.TokenExpConf {
	if c, ok := m.gtcfg[gt]; ok && c != nil {
		return c
	}
	switch gt {
	case config.AuthorizationCode:
		return config.DefaultAuthorizeCodeTokenCfg
	case config.Implicit:
		return config.DefaultImplicitTokenCfg
	case config.PasswordCredentials:
		return config.DefaultPasswordTokenCfg
	case config.ClientCredentials:
		return config.DefaultClientTokenCfg
	}
	return &config.TokenExpConf{}
}

// SetAuthorizeCodeExp set the authorization code expiration time
func (m *ManagerIns) SetAuthorizeCodeExp(exp time.Duration) {
	m.codeExp = exp
}

// SetAuthorizeCodeTokenCfg set the authorization code grant token config
func (m *ManagerIns) SetAuthorizeCodeTokenCfg(cfg *config.TokenExpConf) {
	m.gtcfg[config.AuthorizationCode] = cfg
}

// SetImplicitTokenCfg set the implicit grant token config
func (m *ManagerIns) SetImplicitTokenCfg(cfg *config.TokenExpConf) {
	m.gtcfg[config.Implicit] = cfg
}

// SetPasswordTokenCfg set the password grant token config
func (m *ManagerIns) SetPasswordTokenCfg(cfg *config.TokenExpConf) {
	m.gtcfg[config.PasswordCredentials] = cfg
}

// SetClientTokenCfg set the client grant token config
func (m *ManagerIns) SetClientTokenCfg(cfg *config.TokenExpConf) {
	m.gtcfg[config.ClientCredentials] = cfg
}

// SetRefreshTokenCfg set the refreshing token config
func (m *ManagerIns) SetRefreshTokenCfg(cfg *config.RefreshingConfig) {
	m.rcfg = cfg
}

// MapAuthorizeGenerate mapping the authorize code generate interface
func (m *ManagerIns) MapAuthorizeGenerate(gen generates.AuthorizeGenerate) {
	m.authorizeGenerate = gen
}

// MapAccessGenerate mapping the access token generate interface
func (m *ManagerIns) MapAccessGenerate(gen generates.AccessGenerate) {
	m.accessGenerate = gen
}

// MapClientStorage mapping the client store interface
func (m *ManagerIns) MapClientStorage(stor store.ClientStore) {
	m.clientStore = stor
}

// MustClientStorage mandatory mapping the client store interface
func (m *ManagerIns) MustClientStorage(stor store.ClientStore, err error) {
	if err != nil {
		panic(err.Error())
	}
	m.clientStore = stor
}

// MapTokenStorage mapping the token store interface
func (m *ManagerIns) MapTokenStorage(stor store.TokenStore) {
	m.tokenStore = stor
}

// MustTokenStorage mandatory mapping the token store interface
func (m *ManagerIns) MustTokenStorage(stor store.TokenStore, err error) {
	if err != nil {
		panic(err)
	}
	m.tokenStore = stor
}

// get authorization code data
func (m *ManagerIns) getAuthorizationCode(
	ctx context.Context,
	code string,
) (models.TokenInfo, error) {
	ti, err := m.tokenStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	} else if ti == nil || ti.GetCode() != code || ti.GetCodeCreateAt().Add(ti.GetCodeExpiresIn()).Before(time.Now()) {
		return nil, errors.ErrInvalidAuthorizeCode
	}
	return ti, nil
}

// delete authorization code data
func (m *ManagerIns) delAuthorizationCode(ctx context.Context, code string) error {
	return m.tokenStore.RemoveByCode(ctx, code)
}

// get and delete authorization code data
func (m *ManagerIns) getAndDelAuthorizationCode(
	ctx context.Context,
	tgr *TokenGenerateRequest,
) (models.TokenInfo, error) {
	code := tgr.Code
	ti, err := m.getAuthorizationCode(ctx, code)
	if err != nil {
		return nil, err
	} else if ti.GetClientID() != tgr.ClientID {
		return nil, errors.ErrInvalidAuthorizeCode
	} else if codeURI := ti.GetRedirectURI(); codeURI != "" && codeURI != tgr.RedirectURI {
		return nil, errors.ErrInvalidAuthorizeCode
	}

	err = m.delAuthorizationCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return ti, nil
}

func (m *ManagerIns) validateCodeChallenge(ti models.TokenInfo, ver string) error {
	cc := ti.GetCodeChallenge()
	// early return
	if cc == "" && ver == "" {
		return nil
	}
	if cc == "" {
		return errors.ErrMissingCodeVerifier
	}
	if ver == "" {
		return errors.ErrMissingCodeVerifier
	}
	ccm := ti.GetCodeChallengeMethod()
	if ccm.String() == "" {
		ccm = config.CodeChallengePlain
	}
	if !ccm.Validate(cc, ver) {
		return errors.ErrInvalidCodeChallenge
	}
	return nil
}

// impl interface Manager
// GetClient get the client information
func (m *ManagerIns) GetClient(
	ctx context.Context,
	clientID string,
) (cli models.ClientInfo, err error) {
	cli, err = m.clientStore.GetByID(ctx, clientID)
	if err != nil {
		return
	} else if cli == nil {
		err = errors.ErrInvalidClient
	}
	return
}

// impl interface Manager
// GenerateAuthToken generate the authorization token(code)
func (m *ManagerIns) GenerateAuthToken(
	ctx context.Context,
	rt config.ResponseType,
	tgr *TokenGenerateRequest,
) (models.TokenInfo, error) {

	cli, err := m.GetClient(ctx, tgr.ClientID)
	if err != nil {
		return nil, err
	} else if tgr.RedirectURI != "" {
		if err := m.validateURI(cli.GetDomain(), tgr.RedirectURI); err != nil {
			return nil, err
		}
	}

	ti := models.NewToken()
	ti.SetClientID(tgr.ClientID)
	ti.SetUserID(tgr.UserID)
	ti.SetRedirectURI(tgr.RedirectURI)
	ti.SetScope(tgr.Scope)

	createAt := time.Now()
	td := &generates.GenerateBasic{
		Client:    cli,
		UserID:    tgr.UserID,
		CreateAt:  createAt,
		TokenInfo: ti,
		Request:   tgr.Request,
	}
	switch rt {
	case config.Code:
		codeExp := m.codeExp
		if codeExp == 0 {
			codeExp = config.DefaultCodeExp
		}
		ti.SetCodeCreateAt(createAt)
		ti.SetCodeExpiresIn(codeExp)
		if exp := tgr.AccessTokenExp; exp > 0 {
			ti.SetAccessExpiresIn(exp)
		}
		if tgr.CodeChallenge != "" {
			ti.SetCodeChallenge(tgr.CodeChallenge)
			ti.SetCodeChallengeMethod(tgr.CodeChallengeMethod)
		}

		tv, err := m.authorizeGenerate.Token(ctx, td)
		if err != nil {
			return nil, err
		}
		ti.SetCode(tv)
	case config.Token:
		// set access token expires
		icfg := m.grantConfig(config.Implicit)
		aexp := icfg.AccessTokenExp
		if exp := tgr.AccessTokenExp; exp > 0 {
			aexp = exp
		}
		ti.SetAccessCreateAt(createAt)
		ti.SetAccessExpiresIn(aexp)

		if icfg.IsGenerateRefresh {
			ti.SetRefreshCreateAt(createAt)
			ti.SetRefreshExpiresIn(icfg.RefreshTokenExp)
		}

		tv, rv, err := m.accessGenerate.Token(ctx, td, icfg.IsGenerateRefresh)
		if err != nil {
			return nil, err
		}
		ti.SetAccess(tv)

		if rv != "" {
			ti.SetRefresh(rv)
		}
	}

	err = m.tokenStore.Create(ctx, ti)
	if err != nil {
		return nil, err
	}

	return ti, nil
}

// impl interface Manager
// GenerateAccessToken generate the access token
func (m *ManagerIns) GenerateAccessToken(
	ctx context.Context,
	gt config.GrantType,
	tgr *TokenGenerateRequest,
) (models.TokenInfo, error) {
	cli, err := m.GetClient(ctx, tgr.ClientID)
	if err != nil {
		return nil, err
	}

	if cliPass, ok := cli.(models.ClientPasswordVerifier); ok {
		if !cliPass.VerifyPassword(tgr.ClientSecret) {
			return nil, errors.ErrInvalidClient
		}
	} else if len(cli.GetSecret()) > 0 && tgr.ClientSecret != cli.GetSecret() {
		return nil, errors.ErrInvalidClient
	}
	if tgr.RedirectURI != "" {
		if err := m.validateURI(cli.GetDomain(), tgr.RedirectURI); err != nil {
			return nil, err
		}
	}

	if gt == config.ClientCredentials && cli.IsPublic() {
		return nil, errors.ErrInvalidClient
	}

	if gt == config.AuthorizationCode {
		ti, err := m.getAndDelAuthorizationCode(ctx, tgr)
		if err != nil {
			return nil, err
		}
		if err := m.validateCodeChallenge(ti, tgr.CodeVerifier); err != nil {
			return nil, err
		}
		tgr.UserID = ti.GetUserID()
		tgr.Scope = ti.GetScope()
		if exp := ti.GetAccessExpiresIn(); exp > 0 {
			tgr.AccessTokenExp = exp
		}
	}

	ti := models.NewToken()
	ti.SetClientID(tgr.ClientID)
	ti.SetUserID(tgr.UserID)
	ti.SetRedirectURI(tgr.RedirectURI)
	ti.SetScope(tgr.Scope)

	createAt := time.Now()
	ti.SetAccessCreateAt(createAt)

	// set access token expires
	gcfg := m.grantConfig(gt)
	aexp := gcfg.AccessTokenExp
	if exp := tgr.AccessTokenExp; exp > 0 {
		aexp = exp
	}
	ti.SetAccessExpiresIn(aexp)
	if gcfg.IsGenerateRefresh {
		ti.SetRefreshCreateAt(createAt)
		ti.SetRefreshExpiresIn(gcfg.RefreshTokenExp)
	}

	td := &generates.GenerateBasic{
		Client:    cli,
		UserID:    tgr.UserID,
		CreateAt:  createAt,
		TokenInfo: ti,
		Request:   tgr.Request,
	}

	av, rv, err := m.accessGenerate.Token(ctx, td, gcfg.IsGenerateRefresh)
	if err != nil {
		return nil, err
	}
	ti.SetAccess(av)

	if rv != "" {
		ti.SetRefresh(rv)
	}

	err = m.tokenStore.Create(ctx, ti)
	if err != nil {
		return nil, err
	}

	return ti, nil
}

// impl interface Manager
// RefreshAccessToken refreshing an access token
func (m *ManagerIns) RefreshAccessToken(
	ctx context.Context,
	tgr *TokenGenerateRequest,
) (models.TokenInfo, error) {
	ti, err := m.LoadRefreshToken(ctx, tgr.Refresh)
	if err != nil {
		return nil, err
	}

	cli, err := m.GetClient(ctx, ti.GetClientID())
	if err != nil {
		return nil, err
	}

	oldAccess, oldRefresh := ti.GetAccess(), ti.GetRefresh()

	td := &generates.GenerateBasic{
		Client:    cli,
		UserID:    ti.GetUserID(),
		CreateAt:  time.Now(),
		TokenInfo: ti,
		Request:   tgr.Request,
	}

	rcfg := config.DefaultRefreshTokenCfg
	if v := m.rcfg; v != nil {
		rcfg = v
	}

	ti.SetAccessCreateAt(td.CreateAt)
	if v := rcfg.AccessTokenExp; v > 0 {
		ti.SetAccessExpiresIn(v)
	}

	if v := rcfg.RefreshTokenExp; v > 0 {
		ti.SetRefreshExpiresIn(v)
	}

	if rcfg.IsResetRefreshTime {
		ti.SetRefreshCreateAt(td.CreateAt)
	}

	if scope := tgr.Scope; scope != "" {
		ti.SetScope(scope)
	}

	tv, rv, err := m.accessGenerate.Token(ctx, td, rcfg.IsGenerateRefresh)
	if err != nil {
		return nil, err
	}

	ti.SetAccess(tv)
	if rv != "" {
		ti.SetRefresh(rv)
	}

	if err := m.tokenStore.Create(ctx, ti); err != nil {
		return nil, err
	}

	if rcfg.IsRemoveAccess {
		// remove the old access token
		if err := m.tokenStore.RemoveByAccess(ctx, oldAccess); err != nil {
			return nil, err
		}
	}

	if rcfg.IsRemoveRefreshing && rv != "" {
		// remove the old refresh token
		if err := m.tokenStore.RemoveByRefresh(ctx, oldRefresh); err != nil {
			return nil, err
		}
	}

	if rv == "" {
		ti.SetRefresh("")
		ti.SetRefreshCreateAt(time.Now())
		ti.SetRefreshExpiresIn(0)
	}

	return ti, nil
}

// impl interface Manager
// RemoveAccessToken use the access token to delete the token information
func (m *ManagerIns) RemoveAccessToken(ctx context.Context, access string) error {
	if access == "" {
		return errors.ErrInvalidAccessToken
	}
	return m.tokenStore.RemoveByAccess(ctx, access)
}

// impl interface Manager
// RemoveRefreshToken use the refresh token to delete the token information
func (m *ManagerIns) RemoveRefreshToken(ctx context.Context, refresh string) error {
	if refresh == "" {
		return errors.ErrInvalidAccessToken
	}
	return m.tokenStore.RemoveByRefresh(ctx, refresh)
}

// impl interface Manager
// LoadAccessToken according to the access token for corresponding token information
func (m *ManagerIns) LoadAccessToken(
	ctx context.Context,
	access string,
) (models.TokenInfo, error) {
	if access == "" {
		return nil, errors.ErrInvalidAccessToken
	}

	ct := time.Now()
	ti, err := m.tokenStore.GetByAccess(ctx, access)
	if err != nil {
		return nil, err
	} else if ti == nil || ti.GetAccess() != access {
		return nil, errors.ErrInvalidAccessToken
	} else if ti.GetRefresh() != "" && ti.GetRefreshExpiresIn() != 0 &&
		ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) {
		return nil, errors.ErrExpiredRefreshToken
	} else if ti.GetAccessExpiresIn() != 0 &&
		ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Before(ct) {
		return nil, errors.ErrExpiredAccessToken
	}
	return ti, nil
}

// impl interface Manager
// LoadRefreshToken according to the refresh token for corresponding token information
func (m *ManagerIns) LoadRefreshToken(
	ctx context.Context,
	refresh string,
) (models.TokenInfo, error) {
	if refresh == "" {
		return nil, errors.ErrInvalidRefreshToken
	}

	ti, err := m.tokenStore.GetByRefresh(ctx, refresh)
	if err != nil {
		return nil, err
	} else if ti == nil || ti.GetRefresh() != refresh {
		return nil, errors.ErrInvalidRefreshToken
	} else if ti.GetRefreshExpiresIn() != 0 && // refresh token set to not expire
		ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
		return nil, errors.ErrExpiredRefreshToken
	}
	return ti, nil
}
