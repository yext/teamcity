package teamcity

// Property is a characteristic of a project or build configuration
type Property struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
	Own   bool   `json:"own,omitempty"`
}

// Params is a container for the various properties of a project or build configuration
type Params struct {
	Properties []Property `json:"property,omitempty"`
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
