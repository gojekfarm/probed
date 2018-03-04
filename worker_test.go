package main

import (
	"testing"

	"github.com/rShetty/asyncwait"
	"github.com/stretchr/testify/require"
)

func TestWorkerManagerStartingJobs(t *testing.T) {
	output := 0

	myJobFn := func() {
		output = output + 1
	}

	wm := workerManager{
		workerCount: 2,
		jobFn:       myJobFn,
	}

	wm.start()
	defer wm.stop()

	predicate := func() bool {
		return output == 2
	}

	successful := asyncwait.NewAsyncWait(100, 5).Check(predicate)
	require.True(t, successful)
}
