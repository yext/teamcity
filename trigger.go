package teamcity

import (
	"encoding/json"
	"strconv"
)

// Trigger represents something that kicks off a build type.
type Trigger struct {
	Id                       string
	DependsOn                string
	AfterSuccessfulBuildOnly bool
}

type jsonTrigger struct {
	Id           string        `json:"id,omitempty"`
	Type         string        `json:"type,omitempty"`
	PropertyList *PropertyList `json:"properties,omitempty"`
}

func (t *Trigger) UnmarshalJSON(data []byte) error {
	var jt jsonTrigger
	e := json.Unmarshal(data, &jt)
	if e != nil {
		return e
	}
	*t = Trigger{
		Id:                       jt.Id,
		DependsOn:                jt.PropertyList.Value("dependsOn"),
		AfterSuccessfulBuildOnly: jt.PropertyList.Bool("afterSuccessfulBuildOnly"),
	}
	return nil
}

func (t Trigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonTrigger{
		Id:   t.Id,
		Type: "buildDependencyTrigger",
		PropertyList: NewPropertyList(map[string]string{
			"dependsOn":                t.DependsOn,
			"afterSuccessfulBuildOnly": strconv.FormatBool(t.AfterSuccessfulBuildOnly),
		}),
	})
}
