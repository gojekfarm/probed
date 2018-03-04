package main

type workerManager struct {
	workerCount int
	workQ       chan target
	jobFn       func(chan target)
}

func (wm workerManager) start() {
	for i := 0; i < wm.workerCount; i++ {
		go wm.jobFn(wm.workQ)
	}
}

func (wm workerManager) stop() {
	// Stop go-routines, gracefully ?
}
