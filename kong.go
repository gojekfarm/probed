package main

type kongHealthCheckConfig struct {
	healthCheckPath     string
	healthCheckInterval string
}

type kongHealthCheck struct {
	client            kongClient
	healthCheckConfig kongHealthCheckConfig
}

type kongClient struct {
	kongHost      string
	kongAdminPort string
}

type upstream struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (kc kongClient) upstreams() ([]upstream, error) {
	return nil, nil
}
