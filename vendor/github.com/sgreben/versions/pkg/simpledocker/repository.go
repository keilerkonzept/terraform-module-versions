package simpledocker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Repository represents a Docker repository
type Repository struct {
	URL       string               `json:"URL"`
	V2        *RepositoryV2        `json:"-"`
	DockerHub *RepositoryDockerHub `json:"-"`
}

// Tags returns a list of tags in this repository
func (r *Repository) Tags() (out []struct {
	Name  string
	Image string
}, err error) {
	err = r.Parse()
	if err != nil {
		return nil, err
	}
	switch {
	case r.V2 != nil:
		return r.V2.Tags()
	case r.DockerHub != nil:
		return r.DockerHub.Tags()
	default:
		return nil, errors.New("no docker Repository defined")
	}
}

// RepositoryV2 represents a Docker repository with a V2 API
type RepositoryV2 struct {
	Registry string
	Image    string
}

// Tags returns a list of tags in this repository
func (r *RepositoryV2) Tags() (out []struct {
	Name  string
	Image string
}, err error) {
	url := fmt.Sprintf("https://%s/v2/%s/tags/list", r.Registry, r.Image)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var apiResp struct {
		Tags []string `json:"tags"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	for _, tagName := range apiResp.Tags {
		out = append(out, struct {
			Name  string
			Image string
		}{
			Name:  tagName,
			Image: fmt.Sprintf("%s/%s:%s", r.Registry, r.Image, tagName),
		})
	}
	return out, nil
}

// RepositoryDockerHub represents a Docker Hub repository
type RepositoryDockerHub struct {
	User  string
	Image string
}

// Tags returns a list of tags in this repository
func (r *RepositoryDockerHub) Tags() (out []struct {
	Name  string
	Image string
}, err error) {
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/%s/tags/", r.User, r.Image)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var apiResp struct {
		Count   int     `json:"count"`
		Next    *string `json:"next"`
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	for {
		for _, result := range apiResp.Results {
			tagName := result.Name
			out = append(out, struct {
				Name  string
				Image string
			}{
				Name:  tagName,
				Image: fmt.Sprintf("%s/%s:%s", r.User, r.Image, tagName),
			})
		}
		if apiResp.Next == nil {
			break
		}
		url = *apiResp.Next
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		apiResp.Count = 0
		apiResp.Next = nil
		apiResp.Results = nil
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// Parse parses the repository URL and sets the V2 and DockerHub fields accordingly
func (r *Repository) Parse() error {
	repository := r.URL
	i := strings.IndexByte(repository, '/')
	n := strings.Count(repository, "/")
	switch {
	case n == 0:
		user := "library"
		image := repository
		*r = Repository{
			DockerHub: &RepositoryDockerHub{
				User:  user,
				Image: image,
			},
		}
	case n == 1:
		user := repository[:i]
		image := repository[i+1:]
		*r = Repository{
			DockerHub: &RepositoryDockerHub{
				User:  user,
				Image: image,
			},
		}
	case n >= 2:
		registry, image := repository[:i], repository[i+1:]
		*r = Repository{
			V2: &RepositoryV2{
				Registry: registry,
				Image:    image,
			},
		}
	default:
		return fmt.Errorf("cannot determine registry: %s", repository)
	}
	return nil
}
