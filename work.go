package main

import (
	"errors"
	"fmt"
	"net/http"
)

type pinger struct {
	kongClient KongClient
	pingClient *http.Client
	pingPath   string
	workQ      chan target
}

func (p pinger) start() {
	for t := range p.workQ {
		err := p.pingRequest(t)
		if err != nil {
			p.kongClient.setTargetWeightFor(t.UpstreamID, t.URL, "0")
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
