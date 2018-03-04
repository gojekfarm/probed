package main

type workerManager struct {
	workerCount int
	jobFn       func()
}

func (wm workerManager) start() {
	for i := 0; i < wm.workerCount; i++ {
		go wm.jobFn()
	}
}

func (wm workerManager) stop() {
	// Stop go-routines, gracefully ?
}
