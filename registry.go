package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func registryDiscover(client *http.Client, hostname string) (string, error) {
	var response struct {
		ModulesV1 *string `json:"modules.v1"`
	}
	resp, err := client.Get(fmt.Sprintf("https://%s/.well-known/terraform.json", hostname))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	if response.ModulesV1 == nil {
		return "", fmt.Errorf("no module registry host specified at %q", hostname)
	}
	return *response.ModulesV1, nil
}

func registryListVersions(client *http.Client, baseURL, namespace, name, provider string) ([]string, error) {
	url := fmt.Sprintf("%s%s/%s/%s/versions", baseURL, namespace, name, provider)
	var response struct {
		Modules []struct {
			Versions []struct {
				Version string `json:"version"`
			} `json:"versions"`
		} `json:"modules"`
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %q: %v", url, err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	var versions []string
	for _, m := range response.Modules {
		for _, v := range m.Versions {
			versions = append(versions, v.Version)
		}
	}
	return versions, nil
}
