package versions

import "github.com/sgreben/versions/pkg/simplegit"

// VersionSourceGit describes a version obtained from Git
type VersionSourceGit struct {
	Repository simplegit.Repository
	Reference  string
}

// VersionSourceDocker describes a version obtained from Docker tags
type VersionSourceDocker struct {
	Image string
	Tag   string
}
