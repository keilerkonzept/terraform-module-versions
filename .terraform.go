package main

type dotTerraformModule struct {
	Key     string `json:"Key"`
	Source  string `json:"Source"`
	Version string `json:"Version,omitempty"`
	Dir     string `json:"Dir"`
}

type dotTerraformModules struct {
	Modules []dotTerraformModule `json:"Modules"`
}
