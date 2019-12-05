// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	FakeTsuruTarget  = "https://tsuru.example.com"
	FakeTsuruToken   = "f4k3t0k3n"
	FakeTsuruService = "rpaasv2"
)

func TestNewClient(t *testing.T) {
	t.Run("with default options", func(t *testing.T) {
		c := New("https://rpaas.tsuru.example.com", "admin", "admin")
		assert.Equal(t, &Client{
			address:    "https://rpaas.tsuru.example.com",
			user:       "admin",
			password:   "admin",
			httpClient: &http.Client{Timeout: 30 * time.Second},
		}, c)
	})

	t.Run("with custom timeout", func(t *testing.T) {
		c := NewWithOptions("https://rpaas.tsuru.example.com", "", "", ClientOptions{Timeout: time.Minute})
		assert.Equal(t, &Client{
			address:    "https://rpaas.tsuru.example.com",
			httpClient: &http.Client{Timeout: time.Minute},
		}, c)
	})

	t.Run("client over tsuru w/ default options", func(t *testing.T) {
		c, err := NewClientOverTsuru(FakeTsuruTarget, FakeTsuruToken, FakeTsuruService)
		require.NoError(t, err)
		assert.Equal(t, &Client{
			tsuruTarget:  FakeTsuruTarget,
			tsuruToken:   FakeTsuruToken,
			tsuruService: FakeTsuruService,
			overTsuru:    true,
			httpClient:   &http.Client{Timeout: 30 * time.Second},
		}, c)
	})

	t.Run("client over tsuru w/ custom timeout", func(t *testing.T) {
		c, err := NewClientOverTsuruWithOptions(FakeTsuruTarget, FakeTsuruToken, FakeTsuruService, ClientOptions{Timeout: time.Second})
		require.NoError(t, err)
		assert.Equal(t, &Client{
			tsuruTarget:  FakeTsuruTarget,
			tsuruToken:   FakeTsuruToken,
			tsuruService: FakeTsuruService,
			overTsuru:    true,
			httpClient:   &http.Client{Timeout: time.Second},
		}, c)
	})

	t.Run("client over tsuru getting Tsuru target and token from env vars", func(t *testing.T) {
		require.NoError(t, os.Setenv("TSURU_TARGET", FakeTsuruTarget))
		defer os.Unsetenv("TSURU_TARGET")
		require.NoError(t, os.Setenv("TSURU_TOKEN", FakeTsuruToken))
		defer os.Unsetenv("TSURU_TOKEN")

		c, err := NewClientOverTsuru("", "", FakeTsuruService)
		require.NoError(t, err)
		assert.Equal(t, &Client{
			tsuruTarget:  FakeTsuruTarget,
			tsuruToken:   FakeTsuruToken,
			tsuruService: FakeTsuruService,
			overTsuru:    true,
			httpClient:   &http.Client{Timeout: 30 * time.Second},
		}, c)
	})

	t.Run("client over tsuru without any mandatory argument", func(t *testing.T) {
		var err error
		_, err = NewClientOverTsuru("", FakeTsuruToken, FakeTsuruService)
		assert.EqualError(t, err, "cannot create client without either tsuru target, token or service")
		_, err = NewClientOverTsuru(FakeTsuruTarget, "", FakeTsuruService)
		assert.EqualError(t, err, "cannot create client without either tsuru target, token or service")
		_, err = NewClientOverTsuru(FakeTsuruTarget, FakeTsuruToken, "")
		assert.EqualError(t, err, "cannot create client without either tsuru target, token or service")
	})
}
