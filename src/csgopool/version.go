package csgopool

import (
	"fmt"
)


type Version struct {
	Version int
	Major int
	Minor int
	CodeName string
}

var CurrentVersion *Version

func (v *Version) String() string {
	if len(v.CodeName) > 0 {
		return fmt.Sprintf("%d.%d.%d-%s", v.Version, v.Major, v.Minor, v.CodeName)
	} else {
		return fmt.Sprintf("%d.%d.%d", v.Version, v.Major, v.Minor)
	}
}
