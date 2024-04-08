package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/767829413/advanced-go/api/orgSet"
	"github.com/767829413/advanced-go/api/wxopen"
	"github.com/767829413/advanced-go/http/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	var stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go initHTTP()
	<-stop
	time.Sleep(time.Second)
	fmt.Println("exit successfully")
}

func initHTTP() {
	app := kratos.New(
		kratos.ID("1234"),
		kratos.Name("test"),
		kratos.Version("1.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			NewHTTPServer(),
		),
	)
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func NewHTTPServer() *http.Server {
	srv := http.NewServer(http.Address(":88"))
	// 机构后台配置
	orgSet.RegisterOrgConfHTTPServer(srv, service.NewOrgConfService())
	// 微信配置
	wxopen.RegisterWxopenHTTPServer(srv, service.NewWxopenService())
	return srv
}
