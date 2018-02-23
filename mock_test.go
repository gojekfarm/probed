package main

import "github.com/stretchr/testify/mock"

type mockKongClient struct {
	mock.Mock
}

func (mkc *mockKongClient) upstreams() ([]upstream, error) {
	args := mkc.Called()
	return args.Get(0).([]upstream), args.Error(1)
}

func (mkc *mockKongClient) targetsFor(upstreamID string) ([]target, error) {
	args := mkc.Called(upstreamID)
	return args.Get(0).([]target), args.Error(1)
}

func (mkc *mockKongClient) setTargetWeightFor(upstreamID, targetID, weight string) error {
	args := mkc.Called(upstreamID, targetID, weight)
	return args.Error(0)
}
