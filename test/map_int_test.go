package test

import "testing"

var m1 = map[[12]byte]int{}
var m2 = map[string]int{}

func BenchmarkMapByteInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < 10; i++ {
			var key = [12]byte{}
			copy(key[:], "abc")
			m1[key] = 1
		}
	}
}

func BenchmarkMapInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m2["abc"] = 1
	}
}
