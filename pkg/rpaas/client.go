// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpaas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	proxy      *tsuruServiceProxy
	httpClient *http.Client
}

type ClientArgs struct {
	TsuruTarget  string
	TsuruToken   string
	TsuruService string

	Timeout time.Duration
}

type tsuruServiceProxy struct {
	target  string
	token   string
	service string
}

func NewClient(args ClientArgs) (*Client, error) {
	if args.TsuruTarget == "" || args.TsuruToken == "" || args.TsuruService == "" {
		return nil, fmt.Errorf("cannot create a client without to provide either Tsuru taget, token or service")
	}

	return &Client{
		proxy: &tsuruServiceProxy{
			target:  args.TsuruTarget,
			token:   args.TsuruToken,
			service: args.TsuruService,
		},
		httpClient: &http.Client{
			Timeout: args.Timeout,
		},
	}, nil
}

type Plan struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

func (c *Client) GetPlans(ctx context.Context, instance string) ([]Plan, error) {
	path := fmt.Sprintf("/resources/%s/plans", instance)
	request, err := c.newRequest(http.MethodGet, path, nil, instance)
	if err != nil {
		return nil, err
	}

	body, response, err := c.do(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: Code: %d", response.StatusCode)
	}

	plans := make([]Plan, 0)
	if err = json.Unmarshal(body, &plans); err != nil {
		return nil, fmt.Errorf("unable to decode JSON as plans object: %w", err)
	}

	return plans, nil
}

func (c *Client) newRequest(method, path string, body io.Reader, instance string) (*http.Request, error) {
	request, err := http.NewRequest(method, c.makeURL(path, instance), body)
	if err != nil {
		return nil, fmt.Errorf("unable to build an http request: %w", err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.proxy.token))
	request.Header.Set("Accept", "application/json")
	return request, nil
}

func (c *Client) makeURL(path, instance string) string {
	baseURL := strings.TrimRight(c.proxy.target, "/")
	service := c.proxy.service
	return fmt.Sprintf("%s/services/%s/proxy/%s?callback=%s", baseURL, service, instance, path)
}

func (c *Client) do(ctx context.Context, request *http.Request) ([]byte, *http.Response, error) {
	response, err := c.httpClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, nil, NewClientResponseError(response, err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, NewClientResponseError(response, fmt.Errorf("unable to read body content: %w", err))
	}

	return body, response, nil
}

type ClientResponseError struct {
	Method     string
	URL        string
	StatusCode int

	err      error
	response *http.Response
}

func NewClientResponseError(r *http.Response, err error) *ClientResponseError {
	return &ClientResponseError{
		Method:     r.Request.Method,
		URL:        r.Request.URL.String(),
		StatusCode: r.StatusCode,
		err:        err,
		response:   r,
	}
}

func (e *ClientResponseError) Error() string {
	return fmt.Sprintf("error making HTTP request: URL %s %s - Code: %d", e.Method, e.URL, e.StatusCode)
}

func (e *ClientResponseError) Unwrap() error {
	return e.err
}
