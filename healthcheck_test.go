package main

import (
	"testing"

	"github.com/rShetty/asyncwait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKongHealthCheckStartQueuesTargets(t *testing.T) {
	targetChan := make(chan target, 100)
	mockClient := new(mockKongClient)

	availableUpstreams := []upstream{
		{
			ID:   "1",
			Name: "upstream1",
		},
		{
			ID:   "2",
			Name: "upstream2",
		},
	}

	upstream1Targets := []target{
		{
			ID:     "1.1",
			URL:    "1.2.3.4:80",
			Weight: "1",
		},
		{
			ID:     "1.2",
			URL:    "1.2.3.5:80",
			Weight: "0",
		},
	}

	upstream2Targets := []target{
		{
			ID:     "2.1",
			URL:    "1.2.3.6:80",
			Weight: "100",
		},
		{
			ID:     "2.2",
			URL:    "1.2.3.7:80",
			Weight: "150",
		},
	}

	actualTargets := map[string]target{
		upstream1Targets[0].ID: upstream1Targets[0],
		upstream1Targets[1].ID: upstream1Targets[1],
		upstream2Targets[0].ID: upstream2Targets[0],
		upstream2Targets[1].ID: upstream2Targets[1],
	}

	mockClient.On("upstreams").Return(availableUpstreams, nil)
	mockClient.On("targetsFor", "1").Return(upstream1Targets, nil)
	mockClient.On("targetsFor", "2").Return(upstream2Targets, nil)

	kongHealthCheckConfig := &kongHealthCheckConfig{
		healthCheckPath:     "/ping",
		healthCheckInterval: "10",
	}

	kongHealthCheck := &kongHealthCheck{
		targetChan:            targetChan,
		client:                mockClient,
		kongHealthCheckConfig: kongHealthCheckConfig,
	}

	err := kongHealthCheck.start()
	require.NoError(t, err, "should not have failed to start kong health check")

	predicate := func() bool {
		return len(targetChan) == 4
	}

	successful := asyncwait.NewAsyncWait(100, 10).Check(predicate)
	require.True(t, successful)

	targetMap := make(map[string]target)
	for i := 0; i < 3; i++ {
		target := <-targetChan
		targetMap[target.ID] = target
	}

	for id, target := range targetMap {
		assert.Equal(t, actualTargets[id], target)
	}

	mockClient.AssertExpectations(t)
}
