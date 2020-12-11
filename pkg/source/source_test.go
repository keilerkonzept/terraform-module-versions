package source

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	stringPtr := func(s string) *string { return &s }
	tests := []struct {
		name    string
		raw     string
		want    *Source
		wantErr bool
	}{
		{
			raw: "git@github.com:hashicorp/foo.git",
			want: &Source{
				Git: &Git{
					Remote: "ssh://git@github.com/hashicorp/foo.git",
				},
			},
		},
		{
			raw: "git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=0.10.0",
			want: &Source{
				Git: &Git{
					Remote:   "ssh://git@github.com/keilerkonzept/terraform-module-versions",
					RefValue: stringPtr("0.10.0"),
				},
			},
		},
		{
			raw: "git::git@example.com:foo/bar",
			want: &Source{
				Git: &Git{
					Remote: "ssh://git@example.com/foo/bar",
				},
			},
		},
		{
			raw: "git::git@example.com:foo/bar?ref=0.12.0",
			want: &Source{
				Git: &Git{
					Remote:   "ssh://git@example.com/foo/bar",
					RefValue: stringPtr("0.12.0"),
				},
			},
		},
		{
			raw: "git::git@github.com:keilerkonzept/terraform-module-versions?ref=0.12.0",
			want: &Source{
				Git: &Git{
					Remote:   "ssh://git@github.com/keilerkonzept/terraform-module-versions",
					RefValue: stringPtr("0.12.0"),
				},
			},
		},
		{
			raw: "git::git@github.com:keilerkonzept/terraform-module-versions//pkg/registry?ref=0.12.0",
			want: &Source{
				Git: &Git{
					Remote:     "ssh://git@github.com/keilerkonzept/terraform-module-versions",
					RefValue:   stringPtr("0.12.0"),
					RemotePath: stringPtr("pkg/registry"),
				},
			},
		},
		{
			raw: "hashicorp/consul/aws",
			want: &Source{
				Registry: &Registry{
					Hostname:  "registry.terraform.io",
					Namespace: "hashicorp",
					Name:      "consul",
					Provider:  "aws",
				},
			},
		},
		{
			raw: "example.com:1234/HashiCorp/Consul/aws",
			want: &Source{
				Registry: &Registry{
					Hostname:  "example.com:1234",
					Namespace: "HashiCorp",
					Name:      "Consul",
					Provider:  "aws",
				},
			},
		},
		{
			raw: "github.com/hashicorp/terraform-aws-consul",
			want: &Source{
				Git: &Git{
					Remote: "https://github.com/hashicorp/terraform-aws-consul.git",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}
