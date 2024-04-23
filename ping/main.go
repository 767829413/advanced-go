//go:build linux
// +build linux

package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/767829413/advanced-go/ping/icmp-ping"
	_ "github.com/767829413/advanced-go/ping/mping"
	_ "github.com/767829413/advanced-go/ping/os-ping"
	_ "github.com/767829413/advanced-go/ping/pro-bing"
	etcdclient "go.etcd.io/etcd/client/v3"
)

func main() {
	// osp.OsPing()
	// icmpPing.IcmpPing()
	// proBing.ProBing()
	// mping.Run()
	cli, err := etcdclient.New(etcdclient.Config{
		Endpoints:   []string{"10.20.174.93:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(111, err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = cli.Get(ctx, "test")
	if err != nil {
		fmt.Println(222, err.Error())
	}
	cancel()
}

//
