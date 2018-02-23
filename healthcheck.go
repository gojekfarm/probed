package main

import (
	"strconv"
	"time"
)

type kongHealthCheckConfig struct {
	healthCheckPath     string
	healthCheckInterval string
}

type kongHealthCheck struct {
	targetChan chan target
	client     KongClient
	*kongHealthCheckConfig
}

func (khc *kongHealthCheck) start() error {
	hcInterval, err := strconv.Atoi(khc.healthCheckInterval)
	if err != nil {
		return err
	}

	timeChan := time.NewTicker(time.Millisecond * time.Duration(hcInterval)).C
	go func() {
		select {
		case <-timeChan:
			go khc.monitorHealthOfTargets(khc.targetChan)
		}
	}()

	return nil
}

func (khc *kongHealthCheck) stop() {
	return
}

func (khc *kongHealthCheck) monitorHealthOfTargets(targetChan chan target) {
	upstreams, err := khc.client.upstreams()
	if err != nil {
		//handle
	}

	for _, upstream := range upstreams {
		go khc.fetchAndQueueTargetsFor(upstream.ID, targetChan)
	}

	return
}

func (khc *kongHealthCheck) fetchAndQueueTargetsFor(upstreamID string, targetChan chan target) {
	targets, err := khc.client.targetsFor(upstreamID)
	if err != nil {
		//handle
	}

	for _, target := range targets {
		targetChan <- target
	}
}
