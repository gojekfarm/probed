package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var kongHost = flag.String("kong", "", "kong host")
var kongAdminPort = flag.String("kong-admin-port", "8001", "kong admin port")

var healthCheckInterval = flag.String("health-check-interval", "2000", "healt check interval in ms")
var healthCheckPath = flag.String("health-check-path", "/ping", "path to check for active health check")
var healthCheckType = flag.String("health-check-type", "tcp", "supports http or tcp checks")

var workerCount = flag.Int("worker-count", 100, "no of workers which participate in healthcheck of targets")
var targetsQLen = flag.Int("targets-queue-length", 100, "length of the queue for storing targets")

func main() {
	flag.Parse()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if *kongHost == "" {
		log.Fatalf("`kong` flag did not provide kong host")
	}

	pingQ := make(chan target, *targetsQLen)
	kongClient := newKongClient(*kongHost, *kongAdminPort)

	p := pinger{
		kongClient:      kongClient,
		pingClient:      &http.Client{},
		pingPath:        *healthCheckPath,
		workQ:           pingQ,
		healthCheckType: *healthCheckType,
	}

	wm := newWorkerManager(*workerCount, p.start)

	go wm.start()
	defer wm.stop()

	kongHealthCheckConfig := &kongHealthCheckConfig{
		healthCheckPath:     *healthCheckPath,
		healthCheckInterval: *healthCheckInterval,
	}

	healthCheck, err := newKongHealthCheck(pingQ, kongClient, kongHealthCheckConfig)
	if err != nil {
		log.Fatalf("failed to initialise health checker: %s", err)
	}

	go healthCheck.start()
	defer healthCheck.stop()

	log.Printf("started kong-healthcheck for kong host: %s with interval: %s ms", *kongHost, *healthCheckInterval)
	sig := <-sigChan
	log.Printf("stopping kong-healthcheck, received os signal: %v", sig)
}
