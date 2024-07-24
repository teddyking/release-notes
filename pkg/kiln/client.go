package kiln

import (
	"context"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"gopkg.in/yaml.v2"
)

func New(accessToken string) Client {
	gh, err := factory.NewClient("github", "https://github.com", accessToken)
	if err != nil {
		panic(err)
	}

	return Client{
		gh: gh,
	}
}

type Client struct {
	gh *scm.Client
}

func (c *Client) GetKilnfileAtCommit(commit string, owner string, repo string) (Kilnfile, error) {
	k := Kilnfile{}
	content, _, err := c.gh.Contents.Find(context.Background(), fmt.Sprintf("%s/%s", owner, repo), "Kilnfile", commit)
	if err != nil {
		return k, err
	}

	err = yaml.Unmarshal(content.Data, &k)
	if err != nil {
		return k, err
	}
	return k, nil
}

func (c *Client) GetKilnfileLockAtCommit(commit string, owner string, repo string) (KilnfileLock, error) {
	k := KilnfileLock{}
	content, _, err := c.gh.Contents.Find(context.Background(), fmt.Sprintf("%s/%s", owner, repo), "Kilnfile.lock", commit)
	if err != nil {
		return k, err
	}

	err = yaml.Unmarshal(content.Data, &k)
	if err != nil {
		return k, err
	}
	return k, nil
}
