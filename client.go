package teamcity

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/yext/teamcity/locate"
)

const (
	basePathSuffix  = "/httpAuth/app/rest/"
	projectsPath    = "projects"
	buildsPath      = "builds"
	parametersPath  = "parameters"
	locatorParamKey = "?locator="
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
	if err := c.doRequest("GET", projectsPath, nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectProject gets the project with specified selector
// See https://confluence.jetbrains.com/display/TCD9/REST+API#RESTAPI-ProjectsandBuildConfiguration/TemplatesLists
// for more information about constructing selector.
func (c *Client) SelectProject(selector string) (*Project, error) {
	v := &Project{}
	if err := c.doRequest("GET", path.Join(projectsPath, selector), nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectBuilds gets the build with the specified buildLocator.
// See https://confluence.jetbrains.com/display/TCD9/REST+API#RESTAPI-BuildLocator
// for more information about constructing buildLocator string.
func (c *Client) SelectBuilds(selector string) (*Builds, error) {
	v := &Builds{}
	path := buildsPath + locatorParamKey + selector
	if err := c.doRequest("GET", path, nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// BuildFromId gets the build details for the build with specified id
func (c *Client) BuildFromID(id int) (*Build, error) {
	v := &Build{}
	if err := c.doRequest("GET", path.Join(buildsPath, locate.ById(strconv.Itoa(id)).String()), nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateParameter updates the parameter provided for the specified project name
func (c *Client) UpdateParameter(projectName string, property *Property) (*Property, error) {
	p := path.Join(projectsPath, locate.ByName(projectName).String(), parametersPath, property.Name)
	v := &Property{}
	body, err := json.Marshal(property)
	if err != nil {
		return nil, err
	}
	if err := c.doRequest("PUT", p, body, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Client) doRequest(method string, path string, data []byte, v interface{}) error {
	url := c.host + basePathSuffix + path
	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	rawAuth := []byte(fmt.Sprintf("%v:%v", c.username, c.password))
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(rawAuth))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

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
