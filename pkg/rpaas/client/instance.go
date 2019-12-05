// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tsuru/rpaas-operator/pkg/rpaas/types"
)

type Instance struct {
	Name string
}

func (c *Client) GetPlans(ctx context.Context, inst Instance) ([]types.Plan, error) {
	req, err := c.newRequest("GET", "", "/resources/plans", nil)
	if err != nil {
		return nil, err
	}

	if inst.Name != "" {
		req, err = c.newRequest("GET", inst.Name, fmt.Sprintf("/resources/%s/plans", inst.Name), nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	body, err := getBodyString(resp)
	if err != nil {
		return nil, fmt.Errorf("cannot read the body message: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code \"%d\": body: %v", resp.StatusCode, body)
	}

	var plans []types.Plan
	if err = json.Unmarshal([]byte(body), &plans); err != nil {
		return nil, err
	}

	return plans, nil
}

func (c *Client) GetFlavors(ctx context.Context, inst Instance) ([]types.Flavor, error) {
	req, err := c.newRequest("GET", "", "/resources/flavors", nil)
	if err != nil {
		return nil, err
	}

	if inst.Name != "" {
		req, err = c.newRequest("GET", inst.Name, fmt.Sprintf("/resources/%s/flavors", inst.Name), nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	body, err := getBodyString(resp)
	if err != nil {
		return nil, fmt.Errorf("cannot read the body message: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code \"%d\": body: %v", resp.StatusCode, body)
	}

	var flavors []types.Flavor
	if err = json.Unmarshal([]byte(body), &flavors); err != nil {
		return nil, err
	}

	return flavors, nil
}
