package util

import (
	"hash/fnv"
)

func Str2HashInt(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}
