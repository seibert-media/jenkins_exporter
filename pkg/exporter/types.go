package exporter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BlueJob represents a single job definition from the Jenkins API.
type BlueJob struct {
	Name  string `json:"name"`
	Links struct {
		Self struct {
			URL string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Fail    int `json:"numberOfFailingBranches"`
	Success int `json:"numberOfSuccessfulBranches"`
	Weather int `json:"weatherScore"`
}

// Key generates a usable map key for the job.
func (j *BlueJob) Key() string {
	return j.Name
}

// LegacyJob represents a job from the legacy jenkins api
type LegacyJob struct {
	Name string `json:"name"`
	Jobs []struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"jobs"`
	Health []struct {
		Score int `json:"score"`
	} `json:"healthReport"`
}

// Key generates a usable map key for the job.
func (j *LegacyJob) Key() string {
	return j.Name
}

// BlueRoot represents the root api response from the Jenkins API.
type BlueRoot []BlueJob

// Fetch gathers the root content from the Jenkins API.
func (r *BlueRoot) Fetch(address, username, password string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/blue/rest/search/?q=type:pipeline;organization:jenkins;excludedFromFlattening:jenkins.branch.MultiBranchProject,hudson.matrix.MatrixProject&filter=no-folders&start=0", address),
		nil,
	)

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	res, err := simpleClient().Do(req)

	if err != nil {
		return fmt.Errorf("failed to request root api. %s", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return fmt.Errorf("failed to parse root api. %s", err)
	}

	return nil
}

// LegacyRoot represents the root api response from the Jenkins API.
type LegacyRoot struct {
	Jobs []LegacyJob `json:"jobs"`
}

// Fetch gathers the root content from the Jenkins API.
func (r *LegacyRoot) Fetch(address, username, password string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/json?depth=1&pretty=true&tree=jobs[name,jobs[name,color],healthReport[score]]", address),
		nil,
	)

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	res, err := simpleClient().Do(req)

	if err != nil {
		return fmt.Errorf("failed to request root api. %s", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(r); err != nil {
		return fmt.Errorf("failed to parse root api. %s", err)
	}

	return nil
}
