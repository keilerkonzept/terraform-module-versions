package main

type outputList struct {
	Path              string  `json:"path,omitempty"`
	Name              string  `json:"name,omitempty"`
	Source            string  `json:"source,omitempty"`
	Version           *string `json:"version,omitempty"`
	VersionConstraint *string `json:"versionConstraint,omitempty"`
	Type              string  `json:"type,omitempty"`
}

type outputUpdates struct {
	Path              string   `json:"path,omitempty"`
	Name              string   `json:"name,omitempty"`
	Source            string   `json:"source,omitempty"`
	Version           *string  `json:"version,omitempty"`
	VersionConstraint *string  `json:"versionConstraint,omitempty"`
	Type              string   `json:"type,omitempty"`
	UpdateLatest      string   `json:"latestMatchingUpdate,omitempty"`
	Updates           []string `json:"updates,omitempty"`
	HasMajorUpdate    bool     `json:"hasMajorUpdate,omitempty"`
	HasMinorUpdate    bool     `json:"hasMinorUpdate,omitempty"`
	HasPatchUpdate    bool     `json:"hasPatchUpdate,omitempty"`
}
