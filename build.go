package teamcity

// Builds is a list of builds
type Builds struct {
	Builds []Build `json:"build"`
}

// Build is an instance of a stage in the build chain for a given project
type Build struct {
	Id          int         `json:"id"`
	Number      string      `json:"number"`
	BuildTypeId string      `json:"buildTypeId"`
	BuildType   BuildType   `json:"buildType"`
	Status      string      `json:"status"`
	State       string      `json:"state"`
	Href        string      `json:"href"`
	StatusText  string      `json:"statusText"`
	QueuedDate  string      `json:"queuedDate"`
	StartDate   string      `json:"startDate"`
	FinishDate  string      `json:"finishDate"`
	LastChanges LastChanges `json:"lastChanges"`
}

// BuildType is a type of Build
type BuildType struct {
	Name string `json:"name"`
}

// LastChanges are the list of changes that corresponds to a certain build
type LastChanges struct {
	Changes []Change `json:"change"`
}

// Change is an individual change in a group that corresponds to a certain build
type Change struct {
	Id       int    `json:"id"`
	Version  string `json:"version"`
	Username string `json:"username"`
}
