package main

import (
	_ "github.com/767829413/advanced-go/ping/icmp-ping"
	mping "github.com/767829413/advanced-go/ping/mping"
	_ "github.com/767829413/advanced-go/ping/os-ping"
	_ "github.com/767829413/advanced-go/ping/pro-bing"
)

func main() {
	// osp.OsPing()
	// icmpPing.IcmpPing()
	// proBing.ProBing()
	mping.Run()
}
