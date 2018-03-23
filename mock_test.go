package main

import "github.com/stretchr/testify/mock"

type mockClient struct {
	mock.Mock
}

func (mkc *mockClient) upstreams() ([]upstream, error) {
	args := mkc.Called()
	return args.Get(0).([]upstream), args.Error(1)
}

func (mkc *mockClient) targetsFor(upstreamID string) ([]target, error) {
	args := mkc.Called(upstreamID)
	return args.Get(0).([]target), args.Error(1)
}

func (mkc *mockClient) setTargetWeightFor(upstreamID, targetID string, weight int) error {
	args := mkc.Called(upstreamID, targetID, weight)
	return args.Error(0)
}
