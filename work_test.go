package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rShetty/asyncwait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingSuccess(t *testing.T) {
	mockClient := new(mockKongClient)

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
	pingQ <- target{URL: svr1.URL, Weight: "100", UpstreamID: "upstream1"}
	pingQ <- target{URL: svr2.URL, Weight: "100", UpstreamID: "upstream2"}
	pingQ <- target{URL: svr3.URL, Weight: "100", UpstreamID: "upstream3"}

	mockClient.On("setTargetWeightFor", "upstream3", svr3.URL, "0").Return(nil)

	p := pinger{kongClient: mockClient, pingClient: &http.Client{}, pingPath: *healthCheckPath, workQ: pingQ}
	go p.start()

	predicate := func() bool { return len(pingQ) == 0 }
	successful := asyncwait.NewAsyncWait(100, 5).Check(predicate)
	require.True(t, successful)

	mockClient.AssertExpectations(t)
}
