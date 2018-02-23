package main

type kongHealthCheckConfig struct {
	healthCheckPath     string
	healthCheckInterval string
}

type kongHealthCheck struct {
	client            kongClient
	healthCheckConfig kongHealthCheckConfig
}
