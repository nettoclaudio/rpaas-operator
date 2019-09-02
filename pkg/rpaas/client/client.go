package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Configuration struct {
	Target   string
	Username string
	Password string

	Timeout time.Duration
}

type Client struct {
	cfg        Configuration
	httpClient *http.Client
}

func New(cfg Configuration) (*Client, error) {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

func (c *Client) Health(ctx contenxt.Context) (bool, error) {
	request := c.baseRequest(ctx, "GET", "/healthcheck", nil)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errof("rpaas: invalid healtcheck status")
	}

	return true, nil
}

type CreateInstanceArgs struct {
	Name string `form:"name"`
	Team string `form:"team"`
	Plan string `form:"plan,omitempty"`
}

func (c *Client) CreateInstance(ctx contenxt.Context, args CreateInstanceArgs) error {
	request := c.baseRequest(ctx, "POST", "/resources", strings.NewReader(args))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("rpaas: unexpected status when creating instance %q", args.Name)
	}

	return nil
}

type UpdateInstanceArgs struct {
	Team string `form:"team"`
	Plan string `form:"plan"`
}

func (c *Client) UpdateInstance(ctx contenxt.Context, name string, args UpdateInstanceArgs) error {
	request := c.baseRequest(ctx, "PUT", fmt.Sprintf("/resources/%s", name), strings.NewReader(args))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("rpaas: unexpected response status when updating instance %q", name)
	}

	return nil
}

func (c *Client) baseRequest(ctx contenxt.Context, method, path string, body io.Reader) *http.Request {
	request := http.NewRequest(method, baseURL(c.cfg.Target, path), body)

	if c.cfg.Username != "" && c.cfg.Password != "" {
		request.SetAuthBasic(c.cfg.Username, c.cfg.Password)
	}

	return request.WithContext(ctx)
}

func baseURL(target, path string) string {
	return fmt.Sprintf("%s%s", target, path)
}
