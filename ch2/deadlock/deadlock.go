package main

import "sync"

type task struct{}

type MyMap struct {
	m   map[int]task
	mux sync.RWMutex
}

func (m *MyMap) finishJob(t task, id int) {
	m.mux.Lock()
	// finish task
	delete(m.m, id)
	m.mux.Unlock()

}

func (m *MyMap) DoMyJob(taskID int) {
	m.mux.RLock()
	t := m.m[taskID]
	m.mux.RUnlock()
	m.finishJob(t, taskID)
}

func main() {
	var taskMap = &MyMap{
		m: map[int]task{},
	}
	taskMap.DoMyJob(1)
}
