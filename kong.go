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
	upstreams := []upstream{}

	req, err := http.NewRequest(http.MethodGet, kc.kongAdminURL, nil)
	if err != nil {
		return upstreams, err
	}

	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return upstreams, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return upstreams, err
	}

	upstreamResponse := &upstreamResponse{}

	err = json.Unmarshal(respBytes, upstreamResponse)
	if err != nil {
		return upstreams, err
	}

	return upstreamResponse.Data, nil
}

type target struct {
	ID     string `json:"id"`
	URL    string `json:"target"`
	Weight string `json:"weight"`
}

type targetResponse struct {
	Data []target `json:"data"`
}

func (kc *kongClient) targetsFor(upstreamID string) ([]target, error) {
	targets := []target{}

	req, err := http.NewRequest(http.MethodGet, kc.kongAdminURL, nil)
	if err != nil {
		return targets, err
	}

	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return targets, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return targets, err
	}

	targetResponse := &targetResponse{}

	err = json.Unmarshal(respBytes, targetResponse)
	if err != nil {
		return targets, err
	}

	return targetResponse.Data, nil
}
