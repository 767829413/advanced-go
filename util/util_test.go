package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
)

func TestRandStr(t *testing.T) {
	str := RandStr(66)
	t.Log(str)
	assert.Equal(t, 66, len(str))
}

// go test -timeout 30s -run ^TestConvertAddrs$
func TestConvertAddrs(t *testing.T) {
	ip1 := "8.8.8.8"
	ip2 := "114.114.114.114"

	src := fmt.Sprintf("%s,%s", ip1, ip2)
	addrs, err := ConvertAddrs(src)
	assert.Nil(t, err)
	assert.EqualValues(t, src, strings.Join(addrs, ","))

	src = fmt.Sprintf("%s,%s", ip1, "www.bing.com")
	addrs, err = ConvertAddrs(src)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addrs))

	ip := net.ParseIP(addrs[1])
	assert.NotNil(t, ip)
	t.Log(ip)
	// output: 13.107.21.200

	src = fmt.Sprintf("%s, %s", ip1, ip2)
	addrs, err = ConvertAddrs(src)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(addrs))
	assert.Equal(t, ip1, addrs[0])
	assert.Equal(t, ip2, addrs[1])
	t.Log(addrs[0], addrs[1])

	src = "1.1.0.256, 256.256.256"
	addrs, err = ConvertAddrs(src)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(addrs))
	t.Log(addrs)

	src = "www.ipip.net.xx"
	addrs, err = ConvertAddrs(src)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(addrs))
	t.Log(addrs)

	src = "tools.ipip.net"
	addrs, err = ConvertAddrs(src)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(addrs))
	t.Log(addrs)
}
