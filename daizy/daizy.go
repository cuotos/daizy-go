package daizy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://api-test.daizy.io/api/v1"
)

type API struct {
	organisation string
	authToken    string

	baseURL string

	httpClient *http.Client
	headers    http.Header
}

type Option func(*API) error

// WithBaseURL overrides the default URL for the Daizy API
// default is "https://portal-test.daizy.io/api/v1"
func WithBaseURL(host string) Option {
	return func(api *API) error {
		api.baseURL = host
		return nil
	}
}

// WithTimeout sets the timeout of the http client
func WithTimeout(d time.Duration) Option {
	return func(api *API) error {
		api.httpClient.Timeout = d
		return nil
	}
}

func New(organisation, token string, options ...Option) (*API, error) {

	if organisation == "" {
		return nil, errors.New("organisation ID is required")
	}

	if token == "" {
		return nil, errors.New("authorization token is required")
	}

	c := &API{
		organisation: organisation,
		authToken:    token,
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	if c.baseURL == "" {
		c.baseURL = baseURL
	}

	c.httpClient.Timeout = time.Second * 10

	// Apply each of the optional options
	for _, o := range options {
		o(c)
	}

	return c, nil
}

func (a *API) makeRequest(method, url string, body io.Reader, v interface{}) error {

	fullUrl := fmt.Sprintf("%v%v", a.baseURL, url)

	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.authToken))
	req.Header.Set("Content-Type", fmt.Sprintf("application/json"))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	// TODO: UNTESTED CODE
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status error: %s", resp.Status)
	}

	if v != nil {
		err := json.NewDecoder(resp.Body).Decode(v)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("could not parse json response: %w", err)
		}
	}

	return nil
}
