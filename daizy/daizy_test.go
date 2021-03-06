package daizy

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	//opts = append(nil, opts...)

	client, _ = New(testOrgId, testAuthToken, opts...)

	client.baseURL = server.URL
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
				assert.Equal(t, baseURL, c.baseURL)
				assert.Equal(t, time.Second*10, c.httpClient.Timeout)
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
				WithBaseURL("http://testing-base.com/api/xxx"),
			},
			API{
				organisation: testOrgId,
				authToken:    testAuthToken,
				baseURL:      "http://testing-base.com/api/xxx",
			},
		},
	}

	for _, tc := range tcs {
		c, _ := New(testOrgId, testAuthToken, tc.Options...)

		assert.Equal(t, tc.Expected.baseURL, c.baseURL)
		assert.Equal(t, tc.Expected.organisation, c.organisation)
		assert.Equal(t, tc.Expected.authToken, c.authToken)
	}
}

func TestHeadersAreCorrect(t *testing.T) {
	setup()
	defer teardown()

	called := false
	mux.HandleFunc("/organisation/12345/project/1", func(writer http.ResponseWriter, request *http.Request) {
		called = true
		assert.Equal(t, fmt.Sprintf("Bearer %s", testAuthToken), request.Header.Get("Authorization"))
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(writer, `{
      "name": "aProject",
      "status": "created",
      "user_id": 0,
      "republish_mqtt": true,
      "id": 32,
      "organisation_id": 12
    }`)
	})

	_, err := client.GetProject(1)

	if assert.NoError(t, err) {
		assert.True(t, called, "the endpoint created by the test was not called")
	}
}

//TODO: test that non 200 status return the correct error
func TestNon200Responses(t *testing.T) {
	tcs := []struct {
		ResponseStatus  int
		ResponseMessage string
	}{
		{
			http.StatusBadRequest,
			"A numeric value is required",
		},
	}

	for _, tc := range tcs {
		setup()
		defer teardown()

		mux.HandleFunc("/organisation/12345/project/1", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(tc.ResponseStatus)
			_, _ = fmt.Fprint(writer, `{
  "success": false,
  "errors": [
    {
      "field": "deviceId",
      "type": "NUMERIC",
      "message": "A numeric value is required"
    }
  ]
}`)
		})

		_, err := client.GetProject(1)

		re := &ResponseError{}

		if errors.As(err, &re) {
			assert.Equal(t, "A numeric value is required", re.Error())
			assert.Equal(t, http.StatusBadRequest, re.Status)
		} else {
			t.Fatalf("error was not of type ResponseError: %T", err)
		}
	}
}
