package main

import (
	"context"
	"github.com/767829413/advanced-go/api/orgSet"
	"github.com/767829413/advanced-go/api/wxopen"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func init() {
	InitInternalRpcClient()
	InitWxOpenRpcClient()
}

func main() {
	wxClient := GetWxOpenHTTPClient()
	wxClient.GetAccessToken(context.TODO(), &wxopen.GetAccessTokenRequest{})
	orgClient := GetOrgSetHTTPClient()
	orgClient.GetOrgConf(context.TODO(), &orgSet.GetOrgConfRequest{})
}

var (
	wxopenGrpcConn *http.Client
	httpConn       *http.Client
)

func InitWxOpenRpcClient() {
	//连接http服务
	ctx1, cel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cel()
	var err error
	wxopenGrpcConn, err = http.NewClient(
		ctx1,
		http.WithEndpoint("127.0.0.1:88"),
		http.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
		),
		http.WithTimeout(time.Second*3600),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GetWxOpenConnection() *http.Client {
	return wxopenGrpcConn
}

func GetWxOpenHTTPClient() wxopen.WxopenHTTPClient {
	client := wxopen.NewWxopenHTTPClient(GetWxOpenConnection())
	return client
}

func InitInternalRpcClient() {
	//连接http服务
	ctx1, cel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cel()
	var err error
	httpConn, err = http.NewClient(
		ctx1,
		http.WithEndpoint("127.0.0.1:88"),
		http.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
		),
		http.WithTimeout(time.Second*3600),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConnection() *http.Client {
	return httpConn
}

func GetOrgSetHTTPClient() orgSet.OrgConfHTTPClient {
	client := orgSet.NewOrgConfHTTPClient(GetConnection())
	return client
}
