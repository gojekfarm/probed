package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

const unhealthyNodeWeight = 0
const healthyNodeWeight = 100

type pinger struct {
	kongClient KongClient
	pingClient *http.Client
	pingPath   string
	workQ      chan target
}

func (p pinger) start() {
	for t := range p.workQ {
		log.Printf("pinging target %s", t.URL)
		currentWeight := t.Weight

		err := p.pingRequest(t)
		if err != nil && currentWeight > 0 {
			log.Printf("target %s is down, marking it as unhealthy", t.URL)
			err := p.kongClient.setTargetWeightFor(t.UpstreamID, t.URL, unhealthyNodeWeight)
			if err != nil {
				log.Printf("failed to mark target %s as unhealthy: reason: %s", t.URL, err)
				continue
			}

			continue
		}

		// Previously marked unhealthy node is healthy
		if currentWeight <= 0 && err == nil {
			log.Printf("target %s is up, marking it as healthy", t.URL)
			err := p.kongClient.setTargetWeightFor(t.UpstreamID, t.URL, healthyNodeWeight)
			if err != nil {
				log.Printf("failed to mark target %s as healthy: reason: %s", t.URL, err)
				continue
			}

			continue
		}
	}
}

func (p pinger) pingRequest(t target) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", t.URL, p.pingPath), nil)
	if err != nil {
		return err
	}

	response, err := p.pingClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= http.StatusInternalServerError {
		return errors.New("sever not available")
	}

	return nil
}
