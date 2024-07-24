package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"io"
	"strings"
)

func (c *githubClient) GetPullRequestsForCommits(commits []Commit) ([]PullRequest, error) {
	var fullList []PullRequest

	for _, commit := range commits {
		var pullRequests []PullRequest
		path := fmt.Sprintf("/repos/%s/%s/commits/%s/pulls", c.owner, c.repo, commit.Sha)

		resp, err := c.gh.Do(context.Background(), &scm.Request{
			Method: "GET",
			Path:   path,
		})
		if err != nil {
			return pullRequests, err
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return pullRequests, err
		}

		err = json.Unmarshal(b, &pullRequests)
		if err != nil {
			return pullRequests, err
		}

		if len(pullRequests) > 0 {
			for _, pullRequest := range pullRequests {
				if !containsPullRequest(fullList, pullRequest) {
					fullList = append(fullList, pullRequest)
				}
			}
		} else {
			// this is a commit pushed directly to main - ffs!
			fullList = append(fullList, PullRequest{
				Title:   firstLineOnly(commit.Message),
				HtmlUrl: commit.Url,
			})
		}
	}

	return fullList, nil
}

func firstLineOnly(message string) string {
	parts := strings.Split(message, "\n")
	return parts[0]
}

func containsPullRequest(list []PullRequest, request PullRequest) bool {
	for _, pr := range list {
		if pr.Number == request.Number {
			return true
		}
	}
	return false
}

type PullRequest struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	HtmlUrl string `json:"html_url"`
}
