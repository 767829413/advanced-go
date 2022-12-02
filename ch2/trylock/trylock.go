package trylock

type MyLock struct {
	flag chan struct{}
}

func NewMyLock() *MyLock {
	ml := &MyLock{make(chan struct{}, 1)}
	ml.flag <- struct{}{}
	return ml
}

func (m *MyLock) Lock() {
	<-m.flag
}

func (m *MyLock) UnLock() {
	select {
	case m.flag <- struct{}{}:
	default:
		panic("Unlocked prohibits unlocking")
	}
}

func (m *MyLock) TryLock() bool {
	select {
	case <-m.flag:
		return true
	default:
		return false
	}
}
