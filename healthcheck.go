package main

import (
	"log"
	"strconv"
	"time"
)

type kongHealthCheckConfig struct {
	healthCheckPath     string
	healthCheckInterval string
}

type kongHealthCheck struct {
	ticker     *time.Ticker
	targetChan chan target
	client     KongClient
}

func newKongHealthCheck(targetChan chan target, client KongClient, hcConfig *kongHealthCheckConfig) (*kongHealthCheck, error) {
	hcInterval, err := strconv.Atoi(hcConfig.healthCheckInterval)
	if err != nil {
		return nil, err
	}

	return &kongHealthCheck{
		ticker:     time.NewTicker(time.Millisecond * time.Duration(hcInterval)),
		client:     client,
		targetChan: targetChan,
	}, nil
}

func (khc *kongHealthCheck) start() {
	timeChan := khc.ticker.C

	for range timeChan {
		khc.monitorHealthOfTargets(khc.targetChan)
	}

	return
}

func (khc *kongHealthCheck) stop() {
	khc.ticker.Stop()
}

func (khc *kongHealthCheck) monitorHealthOfTargets(targetChan chan target) {
	upstreams, err := khc.client.upstreams()
	if err != nil {
		log.Printf("failed to fetch upstreams: %s", err)
		return
	}

	for _, upstream := range upstreams {
		go khc.fetchAndQueueTargetsFor(upstream.ID, targetChan)
	}

	return
}

func (khc *kongHealthCheck) fetchAndQueueTargetsFor(upstreamID string, targetChan chan target) {
	targets, err := khc.client.targetsFor(upstreamID)
	if err != nil {
		log.Printf("failed to fetch targets for upstream %s: %s", upstreamID, err)
		return
	}

	for _, target := range targets {
		targetChan <- target
	}
}
