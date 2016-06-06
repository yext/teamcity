package teamcity

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

const (
	basePathSuffix = "/httpAuth/app/rest/"
	projectsPath   = "projects"
)

// Client is an http client and authorization details used to make http requests to TeamCity's API
type Client struct {
	httpClient *http.Client
	host       string
	username   string
	password   string
}

// NewClient creates a new Client with specified authorization details
func NewClient(host, username, password string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		host:       host,
		username:   username,
		password:   password,
	}
}

// ListProjects gets a list of all projects
func (c *Client) ListProjects() (*Projects, error) {
	v := &Projects{}
	if err := c.doRequest("GET", projectsPath, v); err != nil {
		return nil, err
	}
	return v, nil
}

// ProjectByName gets the project with specified selector
// See https://confluence.jetbrains.com/display/TCD9/REST+API#RESTAPI-ProjectsandBuildConfiguration/TemplatesLists
// for more information about constructing selector.
func (c *Client) SelectProject(selector string) (*Project, error) {
	v := &Project{}
	if err := c.doRequest("GET", path.Join(projectsPath, selector), v); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Client) doRequest(method string, path string, v interface{}) error {
	url := c.host + basePathSuffix + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	rawAuth := []byte(fmt.Sprintf("%v:%v", c.username, c.password))
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(rawAuth))
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}
