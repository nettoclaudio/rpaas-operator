// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpaas

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an HTTP client to communicate with RPaaS web API. It can be used
// from two ways: directly or using Tsuru service proxy.
type Client struct {
	args       ClientArgs
	httpClient *http.Client
}

// ClientArgs describes the arguments used to create an instance of Client.
type ClientArgs struct {
	// TsuruTarget stores the Tsuru target.
	TsuruTarget string

	// TsuruToken stores the token used to authenticate into Tsuru target.
	TsuruToken string

	// TsuruService defines the Tsuru service name.
	TsuruService string

	// Timeout defines the limit time for an an HTTP request made by this cliente.
	// Defaults to no timeout.
	Timeout time.Duration
}

func NewClient(args ClientArgs) *Client {
	return &Client{
		args: args,
		httpClient: &http.Client{
			Timeout: args.Timeout,
		},
	}
}

type Plan struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

func (c *Client) GetPlans(ctx context.Context, instance string) ([]Plan, error) {
	endpoint := fmt.Sprintf("/resources/%s/plans", instance)

	request, err := c.newGET(endpoint, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.do(ctx, request)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Client) newGET(endpoint string, values url.Values) (*http.Request, error) {
}

func (c *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	if c.args.TsuruTarget == "" || c.args.TsuruToken == "" || c.args.TsuruService == "" {
		return nil, fmt.Errorf("could not build the tsuru client proxy: missing either target, token or service name")
	}

	baseURL := strings.TrimRight(c.args.TsuruTarget, "/")

	tsuruTarget := c.args.TsuruTarget
}

func (c *Client) do(ctx context.Context, request *http.Request) (*http.Response, error) {
	return c.httpClient.Do(request.WithContext(ctx))
}
