package exporter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Job represents a single job definition from the Jenkins API.
type Job struct {
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
func (j *Job) Key() string {
	return j.Name
}

// Root represents the root api response from the Jenkins API.
type Root []Job

// Fetch gathers the root content from the Jenkins API.
func (r *Root) Fetch(address, username, password string) error {
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
