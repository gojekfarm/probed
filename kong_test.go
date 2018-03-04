package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKongClient(t *testing.T) {
	kclient := newKongClient("127.0.0.1", "9000")

	assert.Equal(t, "127.0.0.1:9000", kclient.kongAdminURL)
	assert.NotNil(t, kclient.httpClient)
}

func TestUpstreamsSuccess(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ "data" : [ {"id": "123-123", "name": "upstream1"}, {"id": "123-124", "name": "upstream2" }] }`))
	}))

	defer httpServer.Close()

	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: httpServer.URL}
	upstreams, err := kclient.upstreams()
	require.NoError(t, err, "should not have failed to get upstreams")

	require.Equal(t, 2, len(upstreams))

	assert.Equal(t, "123-123", upstreams[0].ID)
	assert.Equal(t, "upstream1", upstreams[0].Name)

	assert.Equal(t, "123-124", upstreams[1].ID)
	assert.Equal(t, "upstream2", upstreams[1].Name)
}

func TestUpstreamFailureInvalidURL(t *testing.T) {
	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: "kgp://foo.com"}
	upstreams, err := kclient.upstreams()
	require.Error(t, err, "should have failed to get upstreams")
	assert.Equal(t, 0, len(upstreams))
}

func TestUpstreamsFailureUnmarshal(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ "data" : [ {"id": "123-123", "name": "upstream1"}, {"id": "123-124", "name": "upstream2"] }`))
	}))

	defer httpServer.Close()

	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: httpServer.URL}
	upstreams, err := kclient.upstreams()
	require.Error(t, err, "should have failed to get upstreams")
	require.Equal(t, 0, len(upstreams))
}

func TestTargetsForSuccess(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ "data" : [ {"id": "123-123", "target": "1.2.3.4:8080","weight": "1", "upstream_id": "123-122"}, {"id": "123-124", "target": "1.2.3.5:8080","weight": "0", "upstream_id": "12223"}] }`))
	}))

	defer httpServer.Close()

	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: httpServer.URL}
	targets, err := kclient.targetsFor("upstream1")
	require.NoError(t, err, "should not have failed to get targets")

	require.Equal(t, 2, len(targets))

	assert.Equal(t, "123-123", targets[0].ID)
	assert.Equal(t, "1.2.3.4:8080", targets[0].URL)
	assert.Equal(t, "1", targets[0].Weight)
	assert.Equal(t, "123-122", targets[0].UpstreamID)

	assert.Equal(t, "123-124", targets[1].ID)
	assert.Equal(t, "1.2.3.5:8080", targets[1].URL)
	assert.Equal(t, "0", targets[1].Weight)
	assert.Equal(t, "12223", targets[1].UpstreamID)
}

func TestTargetsForFailureInvalidURL(t *testing.T) {
	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: "kgp://foo.com"}
	upstreams, err := kclient.targetsFor("upstream1")
	require.Error(t, err, "should have failed to get upstreams")
	assert.Equal(t, 0, len(upstreams))
}

func TestTargetsForFailureUnmarshal(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ "data" : [ {"id": "123-123", "target": "1.2.3.4:8080","weight": "1", "upstream_id": "121212"}, {"id": "123-124", "target": "1.2.3.5:8080","weight": "0 }] }`))
	}))

	defer httpServer.Close()

	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: httpServer.URL}
	upstreams, err := kclient.targetsFor("upstream1")
	require.Error(t, err, "should have failed to get upstreams")
	require.Equal(t, 0, len(upstreams))
}

func TestSetTargetWeightForSuccess(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer httpServer.Close()

	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: httpServer.URL}
	err := kclient.setTargetWeightFor("upstream1", "target1", "100")
	require.NoError(t, err, "should not have failed to set target weight")
}

func TestSetTargetWeightForFailure(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer httpServer.Close()

	kclient := &kongClient{httpClient: &http.Client{}, kongAdminURL: httpServer.URL}
	err := kclient.setTargetWeightFor("upstream1", "target1", "100")
	require.Error(t, err, "should have failed to set target weight")
}
