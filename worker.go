package main

import (
	"log"
	"sync"
)

type workerManager struct {
	workerCount int
	jobFn       func()
	doneChan    chan bool

	wg sync.WaitGroup
}

func newWorkerManager(workerCount int, jobFn func()) *workerManager {
	return &workerManager{
		workerCount: workerCount,
		jobFn:       jobFn,
		doneChan:    make(chan bool),
	}
}

func (wm *workerManager) start() {
	for i := 0; i < wm.workerCount; i++ {
		wm.wg.Add(1)
		go func() {
			defer wm.wg.Done()
			wm.jobFn()
		}()
	}
	log.Printf("started %d workers for healthcheck", wm.workerCount)

	<-wm.doneChan
}

func (wm *workerManager) stop() {
	wm.wg.Wait()
	wm.doneChan <- true
}
