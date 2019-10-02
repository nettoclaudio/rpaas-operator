// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpaas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		args         ClientArgs
		wantedError  string
		wantedClient *Client
	}{
		{
			wantedError: "cannot create a client without to provide either Tsuru taget, token or service",
		},
		{
			args: ClientArgs{
				TsuruTarget:  "https://tsuru.test",
				TsuruToken:   "some token",
				TsuruService: "rpaasv2",
				Timeout:      30 * time.Second,
			},
			wantedClient: &Client{
				proxy: &tsuruServiceProxy{
					target:  "https://tsuru.test",
					token:   "some token",
					service: "rpaasv2",
				},
				httpClient: &http.Client{
					Timeout: 30 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		client, err := NewClient(tt.args)
		if tt.wantedError == "" {
			require.NoError(t, err)
		} else {
			require.EqualError(t, err, tt.wantedError)
		}
		assert.Equal(t, client, tt.wantedClient)
	}
}

func TestClient_GetPlans(t *testing.T) {
	tests := []struct {
		instance    string
		handler     http.Handler
		wantedError string
		wantedPlans []Plan
	}{
		{
			instance: "my-instance",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, r.Method, http.MethodGet)
				assert.Equal(t, "Bearer my secure tsuru token", r.Header.Get("Authorization"))
			}),
		},
	}

	for _, tt := range tests {
		server := httptest.NewServer(tt.handler)
		defer server.Close()
		client, err := NewClient(ClientArgs{
			TsuruTarget:  server.URL,
			TsuruToken:   "my secure tsuru token",
			TsuruService: "rpaasv2-test",
		})
		require.NoError(t, err)
		plans, err := client.GetPlans(context.TODO(), tt.instance)
		if tt.wantedError == "" {
			require.NoError(t, err)
		}
		assert.Equal(t, tt.wantedPlans, plans)
	}
}
