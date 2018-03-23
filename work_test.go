package main

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rShetty/asyncwait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingCheckHTTPMarksUnhealthyNodes(t *testing.T) {
	mockClient := new(mockClient)

	svr1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusInternalServerError)
	}))

	pingQ := make(chan target, 10)
	pingQ <- target{URL: svr1.URL, Weight: 100, UpstreamID: "upstream1"}
	pingQ <- target{URL: svr2.URL, Weight: 100, UpstreamID: "upstream2"}
	pingQ <- target{URL: svr3.URL, Weight: 100, UpstreamID: "upstream3"}

	mockClient.On("setTargetWeightFor", "upstream3", svr3.URL, 0).Return(nil)

	p := pinger{client: mockClient, pingClient: &http.Client{}, pingPath: *healthCheckPath, workQ: pingQ, healthCheckType: "http"}
	go p.start()

	predicate := func() bool { return len(pingQ) == 0 }
	successful := asyncwait.NewAsyncWait(100, 5).Check(predicate)
	require.True(t, successful)

	mockClient.AssertExpectations(t)
}

func TestPingCheckHTTPMarksHealthyNodes(t *testing.T) {
	mockClient := new(mockClient)

	svr1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))

	pingQ := make(chan target, 10)
	pingQ <- target{URL: svr1.URL, Weight: 0, UpstreamID: "upstream1"}
	pingQ <- target{URL: svr2.URL, Weight: 100, UpstreamID: "upstream2"}
	pingQ <- target{URL: svr3.URL, Weight: 100, UpstreamID: "upstream3"}

	mockClient.On("setTargetWeightFor", "upstream1", svr1.URL, 100).Return(nil)

	p := pinger{client: mockClient, pingClient: &http.Client{}, pingPath: *healthCheckPath, workQ: pingQ, healthCheckType: "http"}
	go p.start()

	predicate := func() bool { return len(pingQ) == 0 }
	successful := asyncwait.NewAsyncWait(100, 5).Check(predicate)
	require.True(t, successful)

	mockClient.AssertExpectations(t)
}

func TestPingCheckHTTPMarksHealthyNodesFails(t *testing.T) {
	mockClient := new(mockClient)

	svr1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))

	pingQ := make(chan target, 10)
	pingQ <- target{URL: svr1.URL, Weight: 0, UpstreamID: "upstream1"}
	pingQ <- target{URL: svr2.URL, Weight: 0, UpstreamID: "upstream2"}
	pingQ <- target{URL: svr3.URL, Weight: 100, UpstreamID: "upstream3"}

	mockClient.On("setTargetWeightFor", "upstream1", svr1.URL, 100).Return(nil)
	mockClient.On("setTargetWeightFor", "upstream2", svr2.URL, 100).Return(errors.New("failed"))

	p := pinger{client: mockClient, pingClient: &http.Client{}, pingPath: *healthCheckPath, workQ: pingQ, healthCheckType: "http"}
	go p.start()

	predicate := func() bool { return len(pingQ) == 0 }
	successful := asyncwait.NewAsyncWait(100, 5).Check(predicate)
	require.True(t, successful)

	mockClient.AssertExpectations(t)
}

func TestPingCheckHTTPNotMarksNodeAsUnhealthyIfAlreadyUnhealthy(t *testing.T) {
	mockClient := new(mockClient)

	svr1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	svr3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/ping", r.URL.Path)

		w.WriteHeader(http.StatusInternalServerError)
	}))

	pingQ := make(chan target, 10)
	pingQ <- target{URL: svr1.URL, Weight: 100, UpstreamID: "upstream1"}
	pingQ <- target{URL: svr2.URL, Weight: 100, UpstreamID: "upstream2"}
	pingQ <- target{URL: svr3.URL, Weight: 0, UpstreamID: "upstream3"}

	p := pinger{client: mockClient, pingClient: &http.Client{}, pingPath: *healthCheckPath, workQ: pingQ, healthCheckType: "http"}
	go p.start()

	predicate := func() bool { return len(pingQ) == 0 }
	successful := asyncwait.NewAsyncWait(100, 5).Check(predicate)
	require.True(t, successful)

	mockClient.AssertExpectations(t)
}

func TestTCPPortCheckMarksNodesHealthyOrUnhealthy(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	mockClient := new(mockClient)

	pingQ := make(chan target, 4)
	pingQ <- target{URL: "localhost:3000", Weight: 0, UpstreamID: "upstream1"}
	pingQ <- target{URL: "localhost:3000", Weight: 100, UpstreamID: "upstream2"}
	pingQ <- target{URL: "localhost:4000", Weight: 0, UpstreamID: "upstream3"}
	pingQ <- target{URL: "localhost:4000", Weight: 100, UpstreamID: "upstream4"}

	mockClient.On("setTargetWeightFor", "upstream1", "localhost:3000", 100).Return(nil)
	mockClient.On("setTargetWeightFor", "upstream4", "localhost:4000", 0).Return(nil)

	p := pinger{
		client:          mockClient,
		pingClient:      &http.Client{},
		pingPath:        *healthCheckPath,
		workQ:           pingQ,
		healthCheckType: "tcp",
	}
	go p.start()

	predicate := func() bool {
		return len(pingQ) == 0
	}
	successful := asyncwait.NewAsyncWait(100, 5).Check(predicate)
	require.True(t, successful)

	mockClient.AssertExpectations(t)
}
