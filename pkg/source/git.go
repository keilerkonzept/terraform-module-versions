package source

import (
	"fmt"
	"net/url"

	getter "github.com/hashicorp/go-getter"
)

type Git struct {
	Remote     string
	RefValue   *string
	RemotePath *string
}

func parseGitURL(s string) (*Git, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parse git url: %w", err)
	}
	var out Git
	if refValue := u.Query().Get("ref"); refValue != "" {
		out.RefValue = &refValue
		query := u.Query()
		query.Del("ref")
		u.RawQuery = query.Encode()
	}
	out.Remote = u.String()
	if dir, subDir := getter.SourceDirSubdir(out.Remote); subDir != "" {
		out.Remote = dir
		out.RemotePath = &subDir
	}
	return &out, nil
}
