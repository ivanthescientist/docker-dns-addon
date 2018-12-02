package docker

import (
	"strings"
)

func transformContainerName(name string) string {
	return strings.Replace(name, "/", "", -1)
}
