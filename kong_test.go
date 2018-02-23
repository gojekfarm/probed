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

func testUpstreamFailureInvalidURL(t *testing.T) {
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
	require.Error(t, err, "should not have failed to get upstreams")
	require.Equal(t, 0, len(upstreams))
}
