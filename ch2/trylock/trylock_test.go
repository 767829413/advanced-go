package trylock

import (
	"testing"
	"time"
)

func TestMyTryLock(t *testing.T) {
	m := make(map[string]int)
	lk := NewMyLock()
	go func() {
		i := 0
		for {
			if lk.TryLock() {
				m["key"] = i
				lk.UnLock()
				time.Sleep(time.Second)
			} else {
				time.Sleep(time.Second)
			}
			i++
		}
	}()

	go func() {
		for {
			if lk.TryLock() {
				println(m["key"])
				lk.UnLock()
				time.Sleep(time.Second)
			} else {
				time.Sleep(time.Second)
			}
		}
	}()

	timer := time.After(10 * time.Second)
	<-timer
}
