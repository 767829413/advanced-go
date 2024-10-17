package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

func main() {
	uuid := make([]byte, 16)
	makeV7(uuid)
	fmt.Printf("Generated UUID v7: %x\n", uuid)
}

func makeV7(uuid []byte) {
	_ = uuid[15]        // bounds check
	t, s := getV7Time() // 返回毫秒数和序列号

	// 填充前48bit的时间戳
	uuid[0] = byte(t >> 40)
	uuid[1] = byte(t >> 32)
	uuid[2] = byte(t >> 24)
	uuid[3] = byte(t >> 16)
	uuid[4] = byte(t >> 8)
	uuid[5] = byte(t)

	uuid[6] = 0x70 | (0x0F & byte(s>>8)) // 设置版本号7以及后四位存储序列号的前四位
	uuid[7] = byte(s)                    // 存储序列号的后八位

	// 剩余的 uuid[8] ~ uuid[15] 位已填充随机数
	rand.Read(uuid[8:])
}

// 返回毫秒数和序列号
func getV7Time() (uint64, uint16) {
	now := time.Now().UnixMilli()
	seq := make([]byte, 2)
	rand.Read(seq)
	sequence := binary.BigEndian.Uint16(seq)
	return uint64(now), sequence
}
