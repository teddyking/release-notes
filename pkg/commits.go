package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"io"
)

type Change struct {
	Commits []struct {
		Sha     string `json:"sha"`
		HtmlUrl string `json:"html_url"`
		Commit  struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"commits"`
}

type Commit struct {
	Sha     string
	Message string
	Url     string
}

func (c *githubClient) GetCommitsBetween(start string, end string) ([]Commit, error) {
	var commits []Commit

	path := fmt.Sprintf("/repos/%s/%s/compare/%s...%s", c.owner, c.repo, start, end)
	resp, err := c.gh.Do(context.Background(), &scm.Request{
		Method: "GET",
		Path:   path,
	})
	if err != nil {
		return commits, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return commits, err
	}

	change := Change{}
	err = json.Unmarshal(b, &change)
	if err != nil {
		return commits, err
	}

	for _, change := range change.Commits {
		commits = append(commits, Commit{
			Sha:     change.Sha,
			Message: change.Commit.Message,
			Url:     change.HtmlUrl,
		})
	}

	return commits, nil
}
