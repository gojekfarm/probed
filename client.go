package main

// Client is the interface to the Loadbalancer(Kong)
type Client interface {
	upstreams() ([]upstream, error)
	targetsFor(upstreamID string) ([]target, error)
	setTargetWeightFor(upstreamID, targetID string, weight int) error
}
