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

func TestTryLockAfterTime(t *testing.T) {
	lk := NewMyLock()
	c := time.After(20 * time.Second)
	go func() {
		lk.Lock()
		time.Sleep(10 * time.Second)
		lk.UnLock()
	}()

	go func() {
		time.Sleep(1 * time.Second)
		for {
			flag := lk.TryLockAfterTime(3 * time.Second)
			if flag {
				println("lock success")
				lk.UnLock()
				break
			} else {
				println("lock fail,wait !!!")
			}
		}
	}()
	<-c
}
