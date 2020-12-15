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
	baseUri  = "/api/v1"
	baseHost = "https://api-test.daizy.io"
)

type API struct {
	organisation string
	authToken    string

	baseHost string
	baseUri  string

	httpClient *http.Client
	headers    http.Header
}

type Option func(*API) error

// WithBaseHost overrides the default Host for the Daizy API
// default is "https://portal-test.daizy.io"
func WithBaseHost(host string) Option {
	return func(api *API) error {
		api.baseHost = host
		return nil
	}
}

// WithBaseURI overrides the default API URI for the Daizy API
// default is "/api/v1"
func WithBaseURI(uri string) Option {
	return func(api *API) error {
		api.baseUri = uri
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

	// Set this before running the options. If the user intended to set the base uri to "" then it would get set to
	// default due to the comparison looking for empty string ""
	if c.baseUri == "" {
		c.baseUri = baseUri
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	c.httpClient.Timeout = time.Second * 10

	// Apply each of the optional options
	for _, o := range options {
		o(c)
	}

	if c.baseHost == "" {
		c.baseHost = baseHost
	}

	return c, nil
}

func (a *API) makeRequest(method, url string, body io.Reader, v interface{}) error {

	fullUrl := fmt.Sprintf("%v%v%v", a.baseHost, a.baseUri, url)

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

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("could not parse json response: %w", err)
	}

	return nil
}
