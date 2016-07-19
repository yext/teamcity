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

// ByProject gets the Locator for locating by project locator
func ByProject(l Locator) Locator {
	return Locator{"project", fmt.Sprintf("(%v)", l.String())}
}

// BySnapshotDependency gets the Locator for locating by to locator
func BySnapshotDependency(locators ...Locator) Locator {
	var v string
	for _, l := range locators {
		v += l.String() + ","
	}
	return Locator{"snapshotDependency", fmt.Sprintf("(%v)", v[:len(v)-1])}
}

// ByIncludeInitial gets the Locator for locating by includeInitial (used with BySnapshotDependency)
func ByIncludeInitial(b bool) Locator {
	return Locator{"includeInitial", fmt.Sprintf("%v", b)}
}

// ByTo gets the Locator for locating by to locator (used with BySnapshotDependency)
func ByTo(l Locator) Locator {
	return Locator{"to", fmt.Sprintf("(%v)", l.String())}
}
