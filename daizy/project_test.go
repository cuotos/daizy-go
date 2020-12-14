package daizy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCanGetProjects(t *testing.T) {
	setup()
	defer teardown()

	called := false

	mux.HandleFunc("/organisation/12345/projects", func(writer http.ResponseWriter, request *http.Request) {
		called = true
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{
  "projects": [
    {
      "name": "aProject",
      "status": "created",
      "user_id": 0,
      "republish_mqtt": true,
      "id": 32,
      "organisation_id": 12
    }
  ],
  "total": 1,
  "columnFilters": {}
}`))
	})

	expectedProjects := []Project{
		{
			Name:           "aProject",
			Status:         "created",
			UserID:         0,
			RepublishMQTT:  true,
			ID:             32,
			OrganisationID: 12,
		},
	}

	actualProjects, err := client.GetProjects()

	require.True(t, called)

	if assert.NoError(t, err) {
		assert.Equal(t, expectedProjects, actualProjects)
	}
}

func TestCanGetSingleProject(t *testing.T) {
	setup()
	defer teardown()

	called := false

	mux.HandleFunc("/organisation/12345/project/32", func(writer http.ResponseWriter, request *http.Request) {
		called = true
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, `{
      "name": "aProject",
      "status": "created",
      "user_id": 0,
      "republish_mqtt": true,
      "id": 32,
      "organisation_id": 12
    }`)
	})

	expectedProjects := &Project{
		Name:           "aProject",
		Status:         "created",
		UserID:         0,
		RepublishMQTT:  true,
		ID:             32,
		OrganisationID: 12,
	}

	actualProjects, err := client.GetProject(32)

	require.True(t, called, "required endpoint was not called")

	if assert.NoError(t, err) {
		assert.Equal(t, expectedProjects, actualProjects)
	}
}

func TestCreateSingleProject(t *testing.T) {
	setup()
	defer teardown()

	called := false

	mux.HandleFunc("/organisation/12345/project", func(writer http.ResponseWriter, request *http.Request) {
		called = true
		assert.Equal(t, http.MethodPost, request.Method)
		requestBody, _ := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		assert.JSONEq(t, `{
			"name":"aProject",
			"user_id": 444
		}`, string(requestBody))
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(writer, `{
      "name": "aProject",
      "status": "created",
      "user_id": 444,
      "republish_mqtt": true,
      "id": 32,
      "organisation_id": 12
    }`)
	})

	want := &Project{
		Name:           "aProject",
		Status:         "created",
		UserID:         444,
		RepublishMQTT:  true,
		ID:             32,
		OrganisationID: 12,
	}

	create := CreateProject{
		Name:   "aProject",
		UserID: 444,
	}

	actual, err := client.CreateProject(create)

	require.True(t, called, "create project endpoint was not called")

	if assert.NoError(t, err) {
		assert.Equal(t, want, actual)
	}
}
