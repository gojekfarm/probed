package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type kongHealthCheckConfig struct {
	healthCheckPath     string
	healthCheckInterval string
}

type kongHealthCheck struct {
	client            kongClient
	healthCheckConfig kongHealthCheckConfig
}

type kongClient struct {
	httpClient   *http.Client
	kongAdminURL string
}

func newKongClient(kongHost, kongAdminPort string) *kongClient {
	return &kongClient{
		kongAdminURL: fmt.Sprintf("%s:%s", kongHost, kongAdminPort),
		httpClient:   &http.Client{},
	}
}

type upstreamResponse struct {
	Data []upstream `json:"data"`
}

type upstream struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (kc *kongClient) upstreams() ([]upstream, error) {
	req, err := http.NewRequest(http.MethodGet, kc.kongAdminURL, nil)
	if err != nil {
		// handle it
	}

	resp, err := kc.httpClient.Do(req)
	if err != nil {
		// handle it
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle it
	}

	upstreamResponse := &upstreamResponse{}

	err = json.Unmarshal(respBytes, upstreamResponse)
	if err != nil {
		// handle it
	}

	return upstreamResponse.Data, nil
}
