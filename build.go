package teamcity

import (
	"strconv"
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
	Id          int       `json:"id,omitempty"`
	Number      string    `json:"number,omitempty"`
	BuildTypeId string    `json:"buildTypeId,omitempty"`
	BuildType   BuildType `json:"buildType,omitempty"`
	Status      string    `json:"status,omitempty"`
	State       string    `json:"state,omitempty"`
	Href        string    `json:"href,omitempty"`
	StatusText  string    `json:"statusText,omitempty"`
	QueuedDate  Time      `json:"queuedDate,omitempty"`
	StartDate   Time      `json:"startDate,omitempty"`
	FinishDate  Time      `json:"finishDate,omitempty"`
	Changes     Changes   `json:"changes,omitempty"`
	LastChanges Changes   `json:"lastChanges,omitempty"`
	Triggered   Triggered `json:"triggered,omitempty"`
	Comment     Comment   `json:"comment,omitempty"`
	Properties  Params    `json:"properties,omitempty"`
}

// BuildType is a type of Build
type BuildType struct {
	Id                   string                `json:"id,omitempty"`
	Name                 string                `json:"name,omitempty"`
	SnapshotDependencies *SnapshotDependencies `json:"snapshot-dependencies,omitempty"`
	Project              *Project              `json:"project,omitempty"`
}

// BuildTypes is a container for a list of BuildType's
type BuildTypes struct {
	BuildTypes []BuildType `json:"buildType,omitempty"`
}

// Dependency is a build type's artifact or snapshot dependency
type Dependency struct {
	Id              string        `json:"id,omitempty"`
	Type            string        `json:"type,omitempty"`
	SourceBuildType BuildType     `json:"source-buildType,omitempty"`
	PropertyList    *PropertyList `json:"properties,omitempty"`
}

// PropertyList is a list of name-value attributes describing some entity.
type PropertyList struct {
	Count      int        `json:"count"`
	Properties []Property `json:"property"`
}

func NewPropertyList(m map[string]string) *PropertyList {
	var props []Property
	for k, v := range m {
		props = append(props, Property{Name: k, Value: v})
	}
	return &PropertyList{Count: len(props), Properties: props}
}

// Value returns the named property's value, or empty string if not found.
func (pl *PropertyList) Value(name string) string {
	if pl == nil {
		return ""
	}
	for _, v := range pl.Properties {
		if v.Name == name {
			return v.Value
		}
	}
	return ""
}

// Bool returns the named property's boolean value, or false if not found.
func (pl *PropertyList) Bool(name string) bool {
	if pl == nil {
		return false
	}
	var val = pl.Value(name)
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return b
}

// Comment is a description for a Build instance
type Comment struct {
	Text string `json:"text"`
}

// Changes are the list of changes that corresponds to a certain build
type Changes struct {
	Changes []Change `json:"change"`
}

// GetChange returns the most relevant Change describing the build, prioritizing
// Build.Changes over Build.LastChanges out of preference for changes to non-TeamCity repos
func (b *Build) GetChange() Change {
	if len(b.Changes.Changes) > 0 {
		return b.Changes.Changes[0]
	}
	if len(b.LastChanges.Changes) > 0 {
		return b.LastChanges.Changes[0]
	}
	return Change{}
}

// Change is an individual change in a group that corresponds to a certain build
type Change struct {
	Id       int    `json:"id,omitempty"`
	Version  string `json:"version,omitempty"`
	Username string `json:"username,omitempty"`
	Date     Time   `json:"date,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// GetShortVersion returns the first 8 characters of the change version
func (c *Change) GetShortVersion() string {
	var v string
	if len(c.Version) >= 8 {
		v = c.Version[:8]
	}
	return v
}

// BuildsByDate is an interface for sorting a Build array by Date
type BuildsByDate []Build

// Functions for using Golang "sort" package
func (c BuildsByDate) Len() int      { return len(c) }
func (c BuildsByDate) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c BuildsByDate) Less(i, j int) bool {
	timeA := time.Time(c[i].Triggered.Date)
	timeB := time.Time(c[j].Triggered.Date)
	return timeA.Before(timeB)
}

// SnapshotDependencies is a container for SnapshotDependency's
type SnapshotDependencies struct {
	SnapshotDependencies []SnapshotDependency `json:"snapshot-dependency,omitempty"`
}

// SnapshotDependency relates a build type to its source build type
type SnapshotDependency struct {
	SourceBuildType BuildType `json:"source-buildType,omitempty"`
}

// Triggered describes what triggered a particular build
type Triggered struct {
	Date Time `json:"date,omitempty"`
	User User `json:"user,omitempty"`
}

// User describes a user on TeamCity
type User struct {
	Username string `json:"username,omitempty"`
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
