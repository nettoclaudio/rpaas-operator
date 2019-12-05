// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type ClientOptions struct {
	Timeout time.Duration
}

var DefaultClientOptions = ClientOptions{
	Timeout: 30 * time.Second,
}

type Client struct {
	address  string
	user     string
	password string

	tsuruTarget  string
	tsuruToken   string
	tsuruService string
	overTsuru    bool

	httpClient *http.Client
}

func New(address, user, password string) *Client {
	return NewWithOptions(address, user, password, DefaultClientOptions)
}

func NewWithOptions(address, user, password string, opts ClientOptions) *Client {
	return &Client{
		address:    address,
		user:       user,
		password:   password,
		httpClient: newHTTPClient(opts),
	}
}

func NewClientOverTsuru(target, token, service string) (*Client, error) {
	return NewClientOverTsuruWithOptions(target, token, service, DefaultClientOptions)
}

func NewClientOverTsuruWithOptions(target, token, service string, opts ClientOptions) (*Client, error) {
	if target == "" {
		if t, ok := os.LookupEnv("TSURU_TARGET"); ok {
			target = t
		}
	}

	if token == "" {
		if t, ok := os.LookupEnv("TSURU_TOKEN"); ok {
			token = t
		}
	}

	if target == "" || token == "" || service == "" {
		return nil, fmt.Errorf("cannot create client without either tsuru target, token or service")
	}

	return &Client{
		tsuruTarget:  target,
		tsuruToken:   token,
		tsuruService: service,
		overTsuru:    true,
		httpClient:   newHTTPClient(opts),
	}, nil
}

func newHTTPClient(opts ClientOptions) *http.Client {
	return &http.Client{
		Timeout: opts.Timeout,
	}
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req.WithContext(ctx))
}

func (c *Client) newRequest(method, instance, pathName string, body io.Reader) (*http.Request, error) {
	url := c.formatURL(instance, pathName)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if c.tsuruToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.tsuruToken)
	}

	return req, nil
}

func (c *Client) formatURL(instance, pathName string) string {
	if !c.overTsuru {
		return fmt.Sprintf("%s%s", c.address, pathName)
	}

	if instance == "" {
		return fmt.Sprintf("%s/services/proxy/%s?callback=%s", c.tsuruTarget, c.tsuruService, pathName)
	}

	return fmt.Sprintf("%s/services/%s/proxy/%s?callback=%s", c.tsuruTarget, c.tsuruService, instance, pathName)
}

func getBodyString(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body response: %v", err)
	}

	defer resp.Body.Close()
	return string(bodyBytes), nil
}
