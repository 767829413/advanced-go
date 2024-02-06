package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/767829413/advanced-go/open-platform/config"
	"github.com/767829413/advanced-go/open-platform/generates"
	"github.com/767829413/advanced-go/open-platform/manage"
	"github.com/767829413/advanced-go/open-platform/store"
)

var (
	dumpvar   bool
	idvar     string
	secretvar string
	domainvar string
	portvar   int
)

func main() {
	options := &redis.Options{
		Addr:     "redis.rongke-base:6379",
		Password: "",
		DB:       1,
	}
	rdb := redis.NewClient(options)

	manager := manage.NewManager(nil)
	manager.SetAuthorizeCodeExp(config.DefaultCodeExp)
	manager.SetAuthorizeCodeTokenCfg(config.DefaultAuthorizeCodeTokenCfg)
	manager.SetImplicitTokenCfg(config.DefaultImplicitTokenCfg)
	manager.SetPasswordTokenCfg(config.DefaultPasswordTokenCfg)
	manager.SetClientTokenCfg(config.DefaultClientTokenCfg)
	manager.SetRefreshTokenCfg(config.DefaultRefreshTokenCfg)
	manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerateIns())
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	manager.MapTokenStorage(store.NewRedisStoreWithCli(rdb, "rongke-oauth2"))

	cli := &generates.Client{
		ID:     "123456",
		Secret: "QAZWSXEDC",
		Domain: "localhost",
		Public: true,
		UserID: "123456",
	}
	tr := &manage.TokenGenerateRequest{
		ClientID:     cli.ID,
		ClientSecret: cli.Secret,
		UserID:       cli.UserID,
	}

	info, err := manager.GenerateAuthToken(
		context.Background(),
		config.Token,
		tr,
		cli,
	)

	fmt.Println(info, err, info.GetAccess())
}
