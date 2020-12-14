package daizy

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

	// Apply each of the optional options
	for _, o := range options {
		o(c)
	}

	if c.baseHost == "" {
		c.baseHost = baseHost
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c, nil
}

func (a *API) makeRequest(method, url string, body io.Reader) ([]byte, error) {

	fullUrl := fmt.Sprintf("%v%v%v", a.baseHost, a.baseUri, url)

	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.authToken))
	req.Header.Set("Content-Type", fmt.Sprintf("application/json"))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// TODO: UNTESTED CODE
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status error: %s", resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not read the response body: %w", err)
	}

	return respBody, nil
}

func (a *API) internalTestingEndpoint(endpoint string) ([]byte, error) {

	uri := endpoint

	return a.makeRequest(http.MethodGet, uri, nil)
}
