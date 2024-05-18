package monitor

import (
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/gin-gonic/gin"
)

const (
	MonitorPath = "/monitor/notify"
)

var (
	notify = gin.H{
		"code": 0,
		"msg":  "notify success",
	}
	open = make(chan struct{})
)

func NotifyHandler(c *gin.Context) {
	open <- struct{}{}
	c.JSON(http.StatusOK, notify)
}

func RegisterNotifyForProfiling(serverName string, domain string, dir string) {
	go func() {
		var memoryProfile, cpuProfile, traceProfile *os.File
		started := false
		for range open {
			if started {
				pprof.StopCPUProfile()
				trace.Stop()
				pprof.WriteHeapProfile(memoryProfile)
				memoryProfile.Close()
				cpuProfile.Close()
				traceProfile.Close()
				started = false
			} else {
				fileNamePre := fmt.Sprintf("%s/%s_%s", dir, serverName, domain)
				cpuProfile, _ = os.Create(fmt.Sprintf("%s.%s", fileNamePre, "cpu.pprof"))
				memoryProfile, _ = os.Create(fmt.Sprintf("%s.%s", fileNamePre, "memory.pprof"))
				traceProfile, _ = os.Create(fmt.Sprintf("%s.%s", fileNamePre, "runtime.trace"))
				pprof.StartCPUProfile(cpuProfile)
				trace.Start(traceProfile)
				started = true
			}
		}
	}()
}
