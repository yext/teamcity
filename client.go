package teamcity

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/yext/teamcity/locate"
)

var Logger = log.New(ioutil.Discard, "", 0)

const (
	basePathSuffix         = "/httpAuth/app/rest/"
	projectsPath           = "projects"
	buildsPath             = "builds"
	buildTypesPath         = "buildTypes"
	buildQueuePath         = "buildQueue"
	changesPath            = "changes"
	parametersPath         = "parameters"
	templatePath           = "template"
	artifactDependencyPath = "artifact-dependencies"
	snapshotDependencyPath = "snapshot-dependencies"
	triggerPath            = "triggers"
	vcsRootsPath           = "vcs-roots"
	tagsPath               = "tags"

	locatorParamKey = "?locator="

	artifactDependencyType = "artifact_dependency"
	snapshotDependencyType = "snapshot_dependency"

	jsonContentType = "application/json"
	textContentType = "text/plain"
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
	if err := c.doRequest("GET", projectsPath, "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectProject gets the project with specified selector
// See https://confluence.jetbrains.com/display/TCD9/REST+API#RESTAPI-ProjectsandBuildConfiguration/TemplatesLists
// for more information about constructing selector.
func (c *Client) SelectProject(selector string) (*Project, error) {
	v := &Project{}
	if err := c.doRequest("GET", path.Join(projectsPath, selector), "", nil, v); err != nil {
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
	if err := c.doRequest("GET", path, "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// BuildFromId gets the build details for the build with specified id
func (c *Client) BuildFromID(id int) (*Build, error) {
	v := &Build{}
	if err := c.doRequest("GET", path.Join(buildsPath, locate.ById(strconv.Itoa(id)).String()), "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectChange gets the Change with the specified selector
func (c *Client) SelectChange(selector string) (*Change, error) {
	v := &Change{}
	if err := c.doRequest("GET", path.Join(changesPath, selector), "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectBuildType gets the build configuration with the specified selector
func (c *Client) SelectBuildType(selector string) (*BuildType, error) {
	v := &BuildType{}
	if err := c.doRequest("GET", path.Join(buildTypesPath, selector), "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectBuildTypes gets the build configurations with the specified selector
func (c *Client) SelectBuildTypes(selector string) (*BuildTypes, error) {
	v := &BuildTypes{}
	path := buildTypesPath + locatorParamKey + selector
	if err := c.doRequest("GET", path, "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectBuildTypeBuilds gets the builds belonging to the build configuration with the specified selector
func (c *Client) SelectBuildTypeBuilds(selector string) (*Builds, error) {
	v := &Builds{}
	if err := c.doRequest("GET", path.Join(buildTypesPath, selector, buildsPath), "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectVcsRoot gets the VcsRoot belonging to properties specified by the specified selector
func (c *Client) SelectVcsRoot(selector string) (*VcsRoot, error) {
	v := &VcsRoot{}
	if err := c.doRequest("GET", path.Join(vcsRootsPath, selector), "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// TriggerBuildID runs a build for the given build ID and change ID in TeamCity
func (c *Client) TriggerBuildID(buildTypeId string, changeId int, pushDescription string) (*Build, error) {
	v := &Build{}
	build := &Build{
		BuildType: BuildType{
			Id: buildTypeId,
		},
		Properties: Params{
			Properties: []Property{
				Property{
					Name:  "env.PUSH_DESCRIPTION",
					Value: pushDescription,
				},
				Property{
					Name:  "reverse.dep.*.env.PUSH_DESCRIPTION",
					Value: pushDescription,
				},
			},
		},
	}
	if changeId > 0 {
		build.LastChanges = Changes{
			Changes: []Change{
				Change{Id: changeId},
			},
		}
	}
	if len(pushDescription) > 0 {
		build.Comment = Comment{
			Text: pushDescription,
		}
	}
	if err := c.doJSONRequest("POST", buildQueuePath, build, v); err != nil {
		return nil, err
	}
	return v, nil
}

// TriggerBuild runs a build using the given provided *Build.
func (c *Client) TriggerBuild(build *Build, pushDescription string) (*Build, error) {
	if len(pushDescription) > 0 {
		build.Comment = Comment{Text: pushDescription}
	}
	if err := c.doJSONRequest("POST", buildQueuePath, build, build); err != nil {
		return nil, err
	}
	return build, nil
}

// UpdateParameter updates the parameter provided for the specified project name
func (c *Client) UpdateParameter(projectLocator string, property *Property) (*Property, error) {
	p := path.Join(projectsPath, projectLocator, parametersPath, property.Name)
	v := &Property{}
	if err := c.doJSONRequest("PUT", p, property, v); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateBuildTypeParameter updates the parameter provided for the specified build type
func (c *Client) UpdateBuildTypeParameter(buildTypeLocator string, property *Property) (*Property, error) {
	p := path.Join(buildTypesPath, buildTypeLocator, parametersPath, property.Name)
	v := &Property{}
	if err := c.doJSONRequest("PUT", p, property, v); err != nil {
		return nil, err
	}
	return v, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(project *Project) (*Project, error) {
	v := &Project{}
	if err := c.doJSONRequest("POST", projectsPath, project, v); err != nil {
		return nil, err
	}
	return v, nil
}

// CreateBuildType creates a new build type under designated project
func (c *Client) CreateBuildType(projectLocator string, buildType *BuildType) (*BuildType, error) {
	v := &BuildType{}
	p := path.Join(projectsPath, projectLocator, buildTypesPath)
	if err := c.doJSONRequest("POST", p, buildType, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectSnapshotDependency selects a snapshot dependency with given id
func (c *Client) SelectSnapshotDependency(buildTypeSelector string, dependencyId string) (*Dependency, error) {
	v := &Dependency{}
	p := path.Join(buildTypesPath, buildTypeSelector, snapshotDependencyPath, dependencyId)
	if err := c.doRequest("GET", p, "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// SelectArtifactDependencies selects all artifact dependencies for the given build type
func (c *Client) SelectArtifactDependencies(buildTypeSelector string) (*ArtifactDependencies, error) {
	v := &ArtifactDependencies{}
	p := path.Join(buildTypesPath, buildTypeSelector, artifactDependencyPath)
	if err := c.doRequest("GET", p, "", nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

// CreateSnapshotDependency creates a snapshot dependency
func (c *Client) CreateSnapshotDependency(buildTypeSelector string, dependency *Dependency) (*Dependency, error) {
	v := &Dependency{}
	dependency.Type = snapshotDependencyType
	p := path.Join(buildTypesPath, buildTypeSelector, snapshotDependencyPath)
	if err := c.doJSONRequest("POST", p, dependency, v); err != nil {
		return nil, err
	}
	return v, nil
}

// CreateArtifactDependency creates a artifact dependency
func (c *Client) CreateArtifactDependency(buildTypeSelector string, dependency *Dependency) (*Dependency, error) {
	v := &Dependency{}
	dependency.Type = artifactDependencyType
	p := path.Join(buildTypesPath, buildTypeSelector, artifactDependencyPath)
	if err := c.doJSONRequest("POST", p, dependency, v); err != nil {
		return nil, err
	}
	return v, nil
}

// CreateTrigger creates a trigger for a build type
func (c *Client) CreateTrigger(buildTypeSelector string, trigger *Trigger) (*Trigger, error) {
	p := path.Join(buildTypesPath, buildTypeSelector, triggerPath)
	if err := c.doJSONRequest("POST", p, trigger, trigger); err != nil {
		return nil, err
	}
	return trigger, nil
}

// ApplyTemplate applies a build type template to specified build type
func (c *Client) ApplyTemplate(buildTypeSelector string, templateSelector string) (*BuildType, error) {
	v := &BuildType{}
	p := path.Join(buildTypesPath, buildTypeSelector, templatePath)
	if err := c.doRequest("PUT", p, "text/plain", []byte(templateSelector), v); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Client) GetTagByLocator(locator string) (*Tags, error) {
	v := &Tags{}
	p := path.Join(buildsPath, locator, tagsPath)
	if err := c.doJSONRequest("GET", p, nil, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Client) SetTagByLocator(locator string, tags *Tags) (*Tags, error) {
	p := path.Join(buildsPath, locator, tagsPath)
	if err := c.doJSONRequest("PUT", p, tags, tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (c *Client) doJSONRequest(method, path string, t, v interface{}) error {
	body, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := c.doRequest(method, path, jsonContentType, body, v); err != nil {
		return err
	}
	return nil
}

func (c *Client) doRequest(method string, path string, contentType string, data []byte, v interface{}) error {
	Logger.Println(method, path, "\nbody:\n", string(data))
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
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", jsonContentType)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if v != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		Logger.Println("response:\n", string(b))
		if json.Unmarshal(b, v) != nil {
			return errors.New(string(b))
		}
		return nil
	}

	return nil
}
