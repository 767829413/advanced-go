package mping

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/767829413/advanced-go/util"
	"github.com/spf13/pflag"
)

var (
	count        = pflag.IntP("count", "c", 0, "count, 0 means non-setting")
	tos          = pflag.IntP("tos", "z", 0, "tos, 0 means non-setting")
	packetSize   = pflag.IntP("size", "s", 64, "packet size")
	timeout      = pflag.DurationP("timeout", "t", time.Second, "timeout")
	rate         = pflag.IntP("rate", "r", 100, "rate, 100 means 100 packets per second for each target")
	delay        = pflag.IntP("delay", "d", 3, "delay seconds")
	bitflipCheck = pflag.BoolP("bitflip", "", false, "check bitflip")
)

func Run() {
	pflag.ErrHelp = errors.New("")
	pflag.Parse()

	args := pflag.Args()
	if len(args) == 0 {
		pflag.Usage()
		return
	}

	for _, arg := range args {
		tmpArrgs, err := util.ConvertAddrs(arg)
		if err != nil {
			panic(err)
		}
		targetAddrs = append(targetAddrs, tmpArrgs...)
	}

	if *packetSize < len(msgPrefix)+8 {
		*packetSize = len(msgPrefix) + 8
	}

	log.SetFlags(log.Ltime)

	stat = NewBuckets()

	if err := start(); err != nil {
		os.Exit(-1)
	}
}
