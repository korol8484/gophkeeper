package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	cliModel "github.com/korol8484/gophkeeper/internal/client/model"
	"github.com/korol8484/gophkeeper/pkg/model"
	"io"
	"net/http"
	"strings"
	"time"
)

type Render interface {
	Render(map[string]interface{}) error
}

type SecretI interface {
	Format(Render) error
	Title() string
}

// SaveI - hydrate to server model
type SaveI interface {
	ToModel() *model.SecretCreateRequest
}

// Crypto - crypt, decrypt interface
type Crypto interface {
	Encrypt(data []byte, key string) ([]byte, error)
	Decrypt(data []byte, key string) ([]byte, error)
}

// Client - service usecase
type Client struct {
	client *http.Client
	auth   *auth
	cfg    *Config
	crypt  Crypto
}

func NewClient(cfg *Config, crypt Crypto) *Client {
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
		cfg:   cfg,
		crypt: crypt,
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

	request.Header.Add("Content-Type", "application/json")

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

	return errors.New(strings.Trim(string(body), "\n"))
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
	request.Header.Add("Content-Type", "application/json")

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

func (c *Client) Save(ctx context.Context, model SaveI) error {
	var err error

	m := model.ToModel()
	m.Content, err = c.crypt.Encrypt(m.Content, c.auth.passKey)
	if err != nil {
		return err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/user/secret", c.cfg.ServiceHost), bytes.NewReader(b))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", c.auth.token)

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(body))
	}

	return nil
}

func (c *Client) Load(ctx context.Context) ([]cliModel.BaseI, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/user/secret", c.cfg.ServiceHost), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", c.auth.token)

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(body))
	}

	var m []model.Secret
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	for i, v := range m {
		m[i].Content, err = c.crypt.Decrypt(v.Content, c.auth.passKey)
		if err != nil {
			return nil, err
		}
	}

	return cliModel.LoadModels(m), nil
}
