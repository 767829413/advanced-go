package main

import (
	"context"
	"fmt"
	"github.com/767829413/advanced-go/api/orgSet"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"os"
	"time"
)

var (
	// Name is the name of the compiled software.
	Name string = "liveclass"
	// Version is the version of the compiled software.
	Version = "1.0.0"
	id, _   = os.Hostname()
	address = fmt.Sprintf("0.0.0.0:%d", 8888)
)

func main() {
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			NewGRPCServer(),
		),
	)
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func NewGRPCServer() *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
		),
		grpc.Timeout(time.Second * 60),
	}
	opts = append(opts, grpc.Address(address))
	srv := grpc.NewServer(opts...)
	orgSet.RegisterOrgConfServer(srv, HelloService)
	return srv
}

// 定义helloService并实现约定的接口
type helloService struct{ orgSet.UnsafeOrgConfServer }

var HelloService = helloService{}

func (h helloService) GetOrgConf(ctx context.Context, in *orgSet.GetOrgConfRequest) (*orgSet.GetOrgConfResponse, error) {
	resp := new(orgSet.GetOrgConfResponse)
	resp.ShortName = "liveclass"
	return resp, nil
}
func (h helloService) SetOrgConf(ctx context.Context, in *orgSet.SetOrgConfRequest) (*orgSet.SetOrgConfResponse, error) {
	resp := new(orgSet.SetOrgConfResponse)
	resp.Success = "OK"
	return resp, nil
}
