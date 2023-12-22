package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"

	"github.com/767829413/tmp-exec/api/hello"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"os"
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
	hello.RegisterHelloServer(srv, HelloService)
	return srv
}

// 定义helloService并实现约定的接口
type helloService struct{ hello.UnsafeHelloServer }

// HelloService Hello服务
var HelloService = helloService{}

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, in *hello.SayHelloRequest) (*hello.SayHelloResponse, error) {
	resp := new(hello.SayHelloResponse)
	for _, name := range in.Names {
		resp.Hello.Title = fmt.Sprintf("Hello %s.", name)
		resp.Hello.Word = fmt.Sprintf("Gold bless %s.", name)
	}

	return resp, nil
}
