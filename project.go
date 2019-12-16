package teamcity

// Project is an individual project configured in TeamCity
type Project struct {
	Id              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	WebUrl          string   `json:"webUrl,omitempty"`
	Params          Params   `json:"parameters,omitempty"`
	ParentProjectId string   `json:"parentProjectId,omitempty"`
	ParentProject   *Project `json:"parentProject,omitempty"`
}

// Projects is a list of TeamCity projects and aggregate details
type Projects struct {
	Projects []Project `json:"project,omitempty"`
}

// Params is a container for the various properties of a project
type Params struct {
	Properties []Property `json:"property,omitempty"`
}

// Property is a characteristic of a project (e.g. JOB, OWNER, or SERVICE)
type Property struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	Own   bool   `json:"own,omitempty"`
}

// PropertyFromName returns the Property of the given Params with the given target name if it exists
func (params Params) PropertyFromName(target string) Property {
	for _, property := range params.Properties {
		if property.Name == target {
			return property
		}
	}
	return Property{}
}

// PropertyFromName returns the Property of the given Project with the given target name if it exists
func (project Project) PropertyFromName(target string) Property {
	return project.Params.PropertyFromName(target)
}
