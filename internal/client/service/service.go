package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	auth   *auth
	cfg    *Config
}

func NewClient(cfg *Config) *Client {
	def := http.DefaultTransport
	def.(*http.Transport).TLSHandshakeTimeout = 5 * time.Second
	def.(*http.Transport).TLSClientConfig = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
	}

	s := &Client{
		client: &http.Client{
			Transport: def,
		},
		cfg: cfg,
	}

	return s
}

func (c *Client) Auth(ctx context.Context, login, password string) error {
	req := &authRequest{
		Login:    login,
		Password: password,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/user/login", c.cfg.ServiceHost), bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		token := resp.Header.Get("Authorization")
		if len(token) == 0 {
			return errors.New("authorization header not set in response")
		}

		c.auth = newAuthState(token, password)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return errors.New(string(body))
}

func (c *Client) Register(ctx context.Context, login, password string) error {
	req := &authRequest{
		Login:    login,
		Password: password,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/user/register", c.cfg.ServiceHost), bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		token := resp.Header.Get("Authorization")
		if len(token) == 0 {
			return errors.New("authorization header not set in response")
		}

		c.auth = newAuthState(token, password)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return errors.New(string(body))
}
