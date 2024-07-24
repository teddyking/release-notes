package github

import (
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

func New(accessToken string, owner string, repo string) GithubClient {
	gh, err := factory.NewClient("github", "https://github.com", accessToken)
	if err != nil {
		panic(err)
	}

	return &githubClient{
		owner: owner,
		repo:  repo,
		gh:    gh,
	}
}

type GithubClient interface {
	GetCommitsBetween(start string, end string) ([]Commit, error)
	GetPullRequestsForCommits(commits []Commit) ([]PullRequest, error)
	GetLatestTag() (string, error)
}

type githubClient struct {
	owner string
	repo  string
	gh    *scm.Client
}
