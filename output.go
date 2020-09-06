package main

type outputList struct {
	Path              string `json:"path,omitempty"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type,omitempty"`
	Source            string `json:"source,omitempty"`
	VersionConstraint string `json:"constraint,omitempty"`
	Version           string `json:"version,omitempty"`
}

type outputUpdates struct {
	Path                    string `json:"path,omitempty"`
	Name                    string `json:"name,omitempty"`
	Source                  string `json:"source,omitempty"`
	VersionConstraint       string `json:"constraint,omitempty"`
	Version                 string `json:"version,omitempty"`
	ConstraintUpdate        bool   `json:"constraintUpdate,omitempty"`
	Latest                  string `json:"latestMatching,omitempty"`
	MatchingUpdate          bool   `json:"matchingUpdate,omitempty"`
	LatestWithoutConstraint string `json:"latestOverall,omitempty"`
	NonMatchingUpdate       bool   `json:"nonMatchingUpdate,omitempty"`
}
