package main

import (
	"context"
	"fmt"
	"time"

	"github.com/767829413/tmp-exec/api/liveclass"
	fixdiscovery "github.com/767829413/tmp-exec/grpc/fixDiscovery"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	reg "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	grpc2 "google.golang.org/grpc"
)

var (
	grpcConn1 *grpc2.ClientConn

	livecSrvAddr = fmt.Sprintf("127.0.0.1:%d", 9013)
	livecSrv     = &reg.ServiceInstance{
		ID:        "liveclass_service",
		Name:      "liveclass_service",
		Version:   "1.0.0",
		Metadata:  map[string]string{},
		Endpoints: []string{livecSrvAddr},
	}
)

func main() {
	fixdiscovery.InitInternalRpcClientNew(livecSrv, livecSrvAddr)
	conn := GetLiveclassRpcClient()
	req := &liveclass.LiveClassRelatedGroupRequest{}
	req.RelatedGroupId = 103
	resp, err := conn.LiveClassRelatedGroup(context.Background(), req)
	if err != nil {
		log.Error("error: ", err)
		return
	}
	for _, g := range resp.RelatedLiveClass {
		log.Info(g)
	}

}

func GetLiveclassRpcClient() liveclass.LiveClassClient {
	client := liveclass.NewLiveClassClient(fixdiscovery.GetConnectionNew())
	return client
}

func InitInternalRpcClientNew() {
	var err error
	dis := new(livecSrv)
	//连接grpc服务
	ctx1, cel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cel()
	grpcConn1, err = grpc.DialInsecure(
		ctx1,
		grpc.WithEndpoint(livecSrvAddr),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
		),
		grpc.WithTimeout(time.Second*3600),
		grpc.WithDiscovery(dis),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConnectionNew() *grpc2.ClientConn {
	return grpcConn1
}

type fixedDiscovery struct {
	fixSer *reg.ServiceInstance
}

func new(ins *reg.ServiceInstance) *fixedDiscovery {
	return &fixedDiscovery{fixSer: ins}
}

// 根据 serviceName 直接拉取实例列表
func (fdis *fixedDiscovery) GetService(ctx context.Context, serviceName string) ([]*reg.ServiceInstance, error) {
	return []*reg.ServiceInstance{fdis.fixSer}, nil
}

// 根据 serviceName 阻塞式订阅一个服务的实例列表信息
func (fdis *fixedDiscovery) Watch(ctx context.Context, serviceName string) (reg.Watcher, error) {
	w := &watcher{fixSer: fdis.fixSer}
	return w, nil
}

var _ reg.Watcher = (*watcher)(nil)

type watcher struct {
	fixSer *reg.ServiceInstance
}

func (w *watcher) Next() ([]*reg.ServiceInstance, error) {
	return []*reg.ServiceInstance{w.fixSer}, nil
}

func (w *watcher) Stop() error {
	return nil
}
