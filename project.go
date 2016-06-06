package teamcity

// Project is an individual project configured in TeamCity
type Project struct {
	Name   string `json:"name"`
	WebUrl string `json:"webUrl"`
}

// Projects is a list of TeamCity projects and aggregate details
type Projects struct {
	Projects []Project `json:"project"`
}
