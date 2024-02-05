package main

import (
	"context"
	"flag"

	"github.com/go-oauth2/oauth2/v4/generates"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var (
	dumpvar   bool
	idvar     string
	secretvar string
	domainvar string
	portvar   int
)

func init() {
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.StringVar(&idvar, "i", "222222", "The client id being passed in")
	flag.StringVar(&secretvar, "s", "22222222", "The client secret being passed in")
	flag.StringVar(&domainvar, "r", "http://localhost:9094", "The domain of the redirect url")
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {

}

func mainxcc() {
	flag.Parse()
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := store.NewClientStore()
	clientStore.Set(idvar, &models.Client{
		ID:     idvar,
		Secret: secretvar,
		Domain: domainvar,
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetPasswordAuthorizationHandler(
		func(ctx context.Context, clientID, username, password string) (userID string, err error) {
			if username == "test" && password == "test" {
				userID = "test"
			}
			return
		},
	)
}
