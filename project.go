package teamcity

// Project is an individual project configured in TeamCity
type Project struct {
	Name   string `json:"name"`
	WebUrl string `json:"webUrl"`
	Params Params `json:"parameters"`
}

// Projects is a list of TeamCity projects and aggregate details
type Projects struct {
	Projects []Project `json:"project"`
}

// Params is a container for the various properties of a project
type Params struct {
	Properties []Property `json:"property"`
}

// Property is a characteristic of a project (e.g. JOB, OWNER, or SERVICE)
type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Own   bool   `json:"own"`
}
