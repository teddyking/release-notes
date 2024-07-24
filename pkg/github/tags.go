package github

import (
	"context"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
)

func (c *githubClient) GetLatestTag() (string, error) {
	refs, _, err := c.gh.Git.ListTags(context.Background(), fmt.Sprintf("%s/%s", c.owner, c.repo), &scm.ListOptions{})
	if err != nil {
		return "", err
	}

	return refs[0].Name, nil
}
