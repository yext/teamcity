package teamcity

// Project is an individual project configured in TeamCity
type Project struct {
	Id              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	WebUrl          string   `json:"webUrl,omitempty"`
	Params          Params   `json:"parameters,omitempty"`
	ParentProjectId string   `json:"parentProjectId,omitempty"`
	ParentProject   *Project `json:"parentProject,omitempty"`
	Archived        bool     `json:"archived,omitempty"`
}

// Projects is a list of TeamCity projects and aggregate details
type Projects struct {
	Projects []Project `json:"project,omitempty"`
}

// PropertyFromName returns the Property of the given Project with the given target name if it exists
func (project Project) PropertyFromName(target string) Property {
	return project.Params.PropertyFromName(target)
}
