// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ScaleArgs struct {
	Instance
	Quantity int32
}

func (args *ScaleArgs) Validate() error {
	if args.Instance.Name == "" {
		return fmt.Errorf("missing instance name")
	}

	if args.Quantity < int32(0) {
		return fmt.Errorf("replicas number must be greater or equal to zero")
	}

	return nil
}

func (c *Client) Scale(ctx context.Context, args ScaleArgs) error {
	if err := args.Validate(); err != nil {
		return err
	}

	path := fmt.Sprintf("/resources/%s/scale", args.Instance.Name)
	values := url.Values{
		"quantity": []string{fmt.Sprint(args.Quantity)},
	}

	req, err := c.newRequest("POST", args.Instance.Name, path, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	body, err := getBodyString(resp)
	if err != nil {
		return err
	}

	return fmt.Errorf("unexpected status code \"%d\": body: %v", resp.StatusCode, body)
}
