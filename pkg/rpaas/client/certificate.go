// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

type UpdateCertificateArgs struct {
	Instance
	Name        string
	Certificate []byte
	Key         []byte
}

func (args *UpdateCertificateArgs) Validate() error {
	if args.Instance.Name == "" {
		return fmt.Errorf("instance cannot be empty")
	}

	if len(args.Certificate) == 0 {
		return fmt.Errorf("certificate cannot be empty")
	}

	if len(args.Key) == 0 {
		return fmt.Errorf("key cannot be empty")
	}

	return nil
}

func (args *UpdateCertificateArgs) encodeMultipartFormData() (string, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	crtPart, err := writer.CreateFormFile("cert", "tls.crt")
	if err != nil {
		return "", "", err
	}

	if _, err = crtPart.Write(args.Certificate); err != nil {
		return "", "", err
	}

	keyPart, err := writer.CreateFormFile("key", "tls.key")
	if err != nil {
		return "", "", err
	}

	if _, err = keyPart.Write(args.Key); err != nil {
		return "", "", err
	}

	if err = writer.WriteField("name", args.Name); err != nil {
		return "", "", err
	}

	if err = writer.Close(); err != nil {
		return "", "", err
	}

	return body.String(), writer.Boundary(), nil
}

func (c *Client) UpdateCertificate(ctx context.Context, args UpdateCertificateArgs) error {
	if err := args.Validate(); err != nil {
		return err
	}

	raw, boundary, err := args.encodeMultipartFormData()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/resources/%s/certificate", args.Instance.Name)
	req, err := c.newRequest("POST", args.Instance.Name, path, strings.NewReader(raw))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
	resp, err := c.do(ctx, req)

	body, err := getBodyString(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code \"%d\": body: %v", resp.StatusCode, body)
	}

	return nil
}
