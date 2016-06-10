package teamcity

import (
	"strings"
	"time"
)

const (
	dateFormat = "20060102T150405-0700"
)

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
	QueuedDate  Time        `json:"queuedDate"`
	StartDate   Time        `json:"startDate"`
	FinishDate  Time        `json:"finishDate"`
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

// Time is the date in the format TeamCity provides
type Time time.Time

// UnmarshalJSON unmarshals the time using the TeamCity format
func (t *Time) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(dateFormat, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	*t = Time(tt)
	return nil
}

// MarshalJSON marshals the time using the TeamCity format
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(dateFormat) + `"`), nil
}
