package locate

import (
	"fmt"
	"strconv"
)

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
func ById(id int) Locator {
	return Locator{"id", strconv.Itoa(id)}
}

// ByName gets the Locator for locating by name
func ByName(name string) Locator {
	return Locator{"name", name}
}

// ByBuildType gets the Locator for locating by build type locator
func ByBuildType(l Locator) Locator {
	return Locator{"buildType", fmt.Sprintf("(%v)", l.String())}
}
