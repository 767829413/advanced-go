package monitor

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// 添加性能监控
func RegisterHttpForProfiling(path string) {
	/*初始化http服务*/
	port := 82
	address := fmt.Sprintf("0.0.0.0:%d", port)
	router := gin.New()
	pprof.Register(router, path)
	if err := router.Run(address); err != nil {
		fmt.Printf("net pprof is start failed: %s", err.Error())
	}
}

func RegisterRouterForProfiling(router *gin.Engine) {
	// 添加基本的校验机制
	debugGroup := router.Group("/mypprof/:pass", func(c *gin.Context) {
		if c.Param("pass") != "212123453" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})
	pprof.RouteRegister(debugGroup, "debug/pprof")
}
