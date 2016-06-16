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
	Id          int         `json:"id,omitempty"`
	Number      string      `json:"number,omitempty"`
	BuildTypeId string      `json:"buildTypeId,omitempty"`
	BuildType   BuildType   `json:"buildType,omitempty"`
	Status      string      `json:"status,omitempty"`
	State       string      `json:"state,omitempty"`
	Href        string      `json:"href,omitempty"`
	StatusText  string      `json:"statusText,omitempty"`
	QueuedDate  Time        `json:"queuedDate,omitempty"`
	StartDate   Time        `json:"startDate,omitempty"`
	FinishDate  Time        `json:"finishDate,omitempty"`
	LastChanges LastChanges `json:"lastChanges,omitempty"`
	Comment     Comment     `json:"comment,omitempty"`
}

// BuildType is a type of Build
type BuildType struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Comment is a description for a Build instance
type Comment struct {
	Text string `json:"text"`
}

// LastChanges are the list of changes that corresponds to a certain build
type LastChanges struct {
	Changes []Change `json:"change"`
}

// Change is an individual change in a group that corresponds to a certain build
type Change struct {
	Id       int    `json:"id,omitempty"`
	Version  string `json:"version,omitempty"`
	Username string `json:"username,omitempty"`
	Date     Time   `json:"date,omitempty"`
}

// ChangesByDate is an interface for sorting an array of Changes by Date
type ChangesByDate []Change

// Functions for using Golang "sort" package
func (c ChangesByDate) Len() int      { return len(c) }
func (c ChangesByDate) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ChangesByDate) Less(i, j int) bool {
	timeA := time.Time(c[i].Date)
	timeB := time.Time(c[j].Date)
	return timeA.Before(timeB)
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
