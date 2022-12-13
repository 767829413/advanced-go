package timewheel

import (
	"testing"
	"time"
)

func TestTimeWheelAfter(t *testing.T) {
	at := After(10 * time.Second)
	tt := time.After(100 * time.Second)
	go func() {
		i := 0
		for {
			<-at
			println(i)
			i++
			time.Sleep(1 * time.Second)
		}
	}()
	<-tt
}
