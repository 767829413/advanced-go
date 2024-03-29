package proBing

import (
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"log"
	"os"
)

func ProBing() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s host\n", os.Args[0])
		os.Exit(1)
	}
	host := os.Args[1]

	pinger, err := probing.NewPinger(host)
	if err != nil {
		log.Fatal("probing.NewPinger error: ", err)
	}
	// Windows
	pinger.SetPrivileged(true)
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Fatal("pinger.Run error: ", err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	fmt.Printf(
		"reply from %s: seq=%d 往返时延（Round Trip Time）: [MaxRtt=%v,MinRtt=%v,AvgRtt=%v]\n",
		stats.Addr,
		pinger.Count,
		stats.MaxRtt,
		stats.MinRtt,
		stats.AvgRtt,
	)
}
