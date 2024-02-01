package test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"testing"
)

func BenchmarkFmt(b *testing.B) {
	max := new(big.Int).SetInt64(int64(1000000))
	for i := 0; i < b.N; i++ {

		v, _ := rand.Int(rand.Reader, max)
		fmt.Sprint(v.Int64())
	}
}

func BenchmarkStrconv(b *testing.B) {
	max := new(big.Int).SetInt64(int64(1000000))
	for i := 0; i < b.N; i++ {
		v, _ := rand.Int(rand.Reader, max)
		strconv.FormatInt(v.Int64(), 10)
	}
}
