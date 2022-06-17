package utils

import "sync"

type AsyncManager struct {
	wg sync.WaitGroup
}

func (a *AsyncManager) AddTask(funcTask func()) {
	a.wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		funcTask()
	}(&a.wg)
}

func (a *AsyncManager) Wait() {
	a.wg.Wait()
}
