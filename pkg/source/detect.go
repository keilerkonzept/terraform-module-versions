package source

import (
	"fmt"
	"net/url"
	"regexp"

	getter "github.com/hashicorp/go-getter"
)

// https://github.com/hashicorp/terraform/blob/f6d6446701a24c457c14a4a63c113814e3d15144/internal/initwd/getter.go#L23
var goGetterDetectors = []getter.Detector{
	new(getter.GitHubDetector),
	new(getter.GitDetector),
	new(getter.BitBucketDetector),
	new(getter.GCSDetector),
	new(getter.S3Detector),
	new(getter.FileDetector),
}

var forcedGetterRegexp = regexp.MustCompile(`^([A-Za-z0-9]+)::(.+)$`)

func detect(raw string, forceUrl bool) (string, string, error) {
	if (forceUrl) {
		var httpRegexp = regexp.MustCompile(`^([^:]+:\/\/)?([^@]+@)?(.+)$`) // Remove protocol and header indicators
		var colonRegexp = regexp.MustCompile(`^([^?]+):(\D+[\/|?$])`) // remove all colon in miduri (except for pre port indicator)
		raw = httpRegexp.ReplaceAllString(raw, "${3}")
		raw = colonRegexp.ReplaceAllString(raw, "${1}/${2}")
	}
	detected, err := getter.Detect(raw, ".", goGetterDetectors)
	if err != nil {
		return "", "", fmt.Errorf("detect source type: %w", err)
	}
	var forced string
	if match := forcedGetterRegexp.FindStringSubmatch(detected); match != nil {
		forced = match[1]
		detected = match[2]
		return forced, detected, nil
	}
	if u, _ := url.Parse(detected); u != nil {
		forced = u.Scheme
	}
	return forced, detected, nil
}
