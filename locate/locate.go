package locate

import "fmt"

// Locator is a key, value used to locate various TeamCity entities
type Locator struct {
	key   string
	value string
}

// String converts the locator to a string in the form key:value
func (l Locator) String() string {
	return l.key + ":" + l.value
}

// ById gets the Locator for locating by id
func ById(id string) Locator {
	return Locator{"id", id}
}

// ByName gets the Locator for locating by name
func ByName(name string) Locator {
	return Locator{"name", name}
}

// ByVersion gets the Locator for locating a Change by version
func ByVersion(version string) Locator {
	return Locator{"version", version}
}

// ByBuildType gets the Locator for locating by build type locator
func ByBuildType(l Locator) Locator {
	return Locator{"buildType", fmt.Sprintf("(%v)", l.String())}
}

// ByAffectedProject gets the Locator for locating by affected project locator
func ByAffectedProject(l Locator) Locator {
	return Locator{"affectedProject", fmt.Sprintf("(%v)", l.String())}
}
