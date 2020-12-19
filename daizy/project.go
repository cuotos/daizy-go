package daizy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	UserID         int    `json:"user_id"`
	RepublishMQTT  bool   `json:"republish_mqtt"`
	ID             int    `json:"id"`
	OrganisationID int    `json:"organisation_id"`
}

type ProjectListResponse struct {
	Projects      []Project   `json:"projects"`
	Total         int         `json:"total"`
	ColumnFilters interface{} `json:"columnFilters,omitempty"`
}

type CreateProject struct {
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

type UpdateProject CreateProject

// GetProjects returns a slice of all the Projects for the Organisation
func (a *API) GetProjects() ([]Project, error) {
	uri := fmt.Sprintf("/organisation/%s/projects", a.organisation)

	resp, err := a.makeRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	projectsResp := &ProjectListResponse{}

	err = json.Unmarshal(resp, projectsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateProjectResponse: %w", err)
	}

	return projectsResp.Projects, nil
}

// GetProject return a single Project
func (a *API) GetProject(id int) (*Project, error) {
	//TODO extract the org/<id> bit out to the client as its always the same
	uri := fmt.Sprintf("/organisation/%s/project/%d", a.organisation, id)

	resp, err := a.makeRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	project := &Project{}

	err = json.Unmarshal(resp, project)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateProjectResponse: %w", err)
	}

	return project, nil
}

// CreateProject creates a project
func (a *API) CreateProject(project CreateProject) (*Project, error) {

	uri := fmt.Sprintf("/organisation/%s/project", a.organisation)

	createProjectBytes, err := json.Marshal(project)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall project to json: %w", err)
	}

	resp, err := a.makeRequest(http.MethodPost, uri, bytes.NewReader(createProjectBytes))

	createProjectResponse := &Project{}

	err = json.Unmarshal(resp, createProjectResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateProjectResponse: %w", err)
	}

	return createProjectResponse, nil
}

// DeleteProject deleted project with the provided ID
func (a *API) DeleteProject(id int) error {
	uri := fmt.Sprintf("/organisation/%s/project/%d", a.organisation, id)

	_, err := a.makeRequest(http.MethodDelete, uri, nil)

	return err
}

// UpdateProject updates a Project with the given ID
func (a *API) UpdateProject(id int, project UpdateProject) (*Project, error) {
	uri := fmt.Sprintf("/organisation/%s/project/%d", a.organisation, id)

	createProjectBytes, err := json.Marshal(project)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall project to json: %w", err)
	}

	resp, err := a.makeRequest(http.MethodPut, uri, bytes.NewReader(createProjectBytes))
	if err != nil {
		return nil, fmt.Errorf("unable to marshall project to json: %w", err)
	}

	updateProjectResponse := &Project{}

	err = json.Unmarshal(resp, updateProjectResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CreateProjectResponse: %w", err)
	}

	return updateProjectResponse, nil
}
