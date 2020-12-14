package daizy

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testOrgId     = "12345"
	testAuthToken = "testtoken"
)

var (
	mux    *http.ServeMux
	client *API
	server *httptest.Server
)

func setup(opts ...Option) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// prepend so any provided values will override
	opts = append([]Option{WithBaseURI("")}, opts...)

	client, _ = New(testOrgId, testAuthToken, opts...)

	client.baseHost = server.URL
}

func teardown() {
	server.Close()
}

func TestCanCreateNewClient(t *testing.T) {

	tcs := []struct {
		InputOrg   string
		InputToken string
		Error      error
	}{
		{
			testOrgId,
			testAuthToken,
			nil,
		},
		{
			"",
			testAuthToken,
			errors.New("organisation ID is required"),
		},
		{
			testOrgId,
			"",
			errors.New("authorization token is required"),
		},
	}

	for _, tc := range tcs {

		c, err := New(tc.InputOrg, tc.InputToken)

		if tc.Error == nil {
			if assert.NoError(t, err, "valid client args should not error") {
				assert.Equal(t, tc.InputOrg, c.organisation)
				assert.Equal(t, tc.InputToken, c.authToken)
				assert.Equal(t, baseHost, c.baseHost)
				assert.Equal(t, baseUri, c.baseUri)
			}
		} else {
			if assert.Error(t, err) {
				assert.EqualError(t, err, tc.Error.Error())
			}
		}
	}
}

func TestCanCreateClientWithCustomOptions(t *testing.T) {

	tcs := []struct {
		Options  []Option
		Expected API
	}{
		{
			[]Option{
				WithBaseHost("http://testing-base.com"),
				WithBaseURI("/api/xxx"),
			},
			API{
				organisation: testOrgId,
				authToken:    testAuthToken,
				baseHost:     "http://testing-base.com",
				baseUri:      "/api/xxx",
			},
		},
	}

	for _, tc := range tcs {
		c, _ := New(testOrgId, testAuthToken, tc.Options...)

		assert.Equal(t, tc.Expected.baseUri, c.baseUri)
		assert.Equal(t, tc.Expected.baseHost, c.baseHost)
		assert.Equal(t, tc.Expected.organisation, c.organisation)
		assert.Equal(t, tc.Expected.authToken, c.authToken)
	}
}

func TestHeadersAreCorrect(t *testing.T) {
	setup(WithBaseURI("/new/base/uri"))
	defer teardown()

	called := false
	mux.HandleFunc("/new/base/uri/randomEndpointForTestingOnly", func(writer http.ResponseWriter, request *http.Request) {
		called = true
		assert.Equal(t, fmt.Sprintf("Bearer %s", testAuthToken), request.Header.Get("Authorization"))
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
	})

	_, err := client.internalTestingEndpoint("/randomEndpointForTestingOnly")

	if assert.NoError(t, err) {
		assert.True(t, called, "the endpoint created by the test was not called")
	}
}

//TODO: test that non 200 status return the correct error
