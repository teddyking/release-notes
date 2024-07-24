package main

import (
	"fmt"
	"github.com/garethjevans/release-notes/pkg/config"
	"github.com/garethjevans/release-notes/pkg/github"
	"github.com/garethjevans/release-notes/pkg/kiln"
	"net/url"
	"os"
	"strings"
)

func main() {
	releaseNotesConfig := config.MustLoad(".github/release-notes.yml")
	accessToken := os.Getenv("GITHUB_TOKEN")
	if accessToken == "" {
		panic("no access token is available via GITHUB_TOKEN")
	}

	currentRelease := "HEAD"

	c := github.New(accessToken, releaseNotesConfig.Owner, releaseNotesConfig.Repo)

	previousRelease, err := c.GetLatestTag()
	if err != nil {
		panic(err)
	}

	commits, err := c.GetCommitsBetween(previousRelease, currentRelease)
	if err != nil {
		panic(err)
	}

	prs, err := c.GetPullRequestsForCommits(commits)
	if err != nil {
		panic(err)
	}

	fmt.Printf("# %s\n", releaseNotesConfig.Title)
	for _, pr := range prs {
		if pr.Number == 0 {
			fmt.Printf("* [%s](%s)\n", pr.Title, pr.HtmlUrl)
		} else {
			fmt.Printf("* #%d [%s](%s)\n", pr.Number, pr.Title, pr.HtmlUrl)
		}
	}

	for _, include := range releaseNotesConfig.Includes {
		if include.Type == "kiln" {
			k := kiln.New(accessToken)

			baseKilnfile, err := k.GetKilnfileAtCommit(previousRelease, releaseNotesConfig.Owner, releaseNotesConfig.Repo)
			if err != nil {
				panic(err)
			}

			baseKilnfileLock, err := k.GetKilnfileLockAtCommit(previousRelease, releaseNotesConfig.Owner, releaseNotesConfig.Repo)
			if err != nil {
				panic(err)
			}

			currentKilnfileLock, err := k.GetKilnfileLockAtCommit(currentRelease, releaseNotesConfig.Owner, releaseNotesConfig.Repo)
			if err != nil {
				panic(err)
			}

			genaiBase := baseKilnfileLock.GetVersionForRelease(include.Name)
			genaiCurrent := currentKilnfileLock.GetVersionForRelease(include.Name)

			if genaiBase != genaiCurrent {
				fmt.Println("\n\n# Dependency Change")
				fmt.Println("| Dependency | Type | From | To |")
				fmt.Println("| ---------- | ---- | ---- | -- |")
				fmt.Printf("| %s | %s | %s | %s |\n\n", include.Name, include.Type, genaiBase, genaiCurrent)
			}

			gitRepo := baseKilnfile.GetGithubRepositoryForRelease(include.Name)
			o, r := MustExtractOwnerAndRepoFromGitUrl(gitRepo)

			c = github.New(accessToken, o, r)

			commits, err = c.GetCommitsBetween(
				fmt.Sprintf("v%s", genaiBase),
				fmt.Sprintf("v%s", genaiCurrent))
			if err != nil {
				panic(err)
			}

			prs, err = c.GetPullRequestsForCommits(commits)
			if err != nil {
				panic(err)
			}

			fmt.Printf("\n\n# %s\n", include.Title)
			for _, pr := range prs {
				if pr.Number == 0 {
					fmt.Printf("* [%s](%s)\n", pr.Title, pr.HtmlUrl)
				} else {
					fmt.Printf("* #%d [%s](%s)\n", pr.Number, pr.Title, pr.HtmlUrl)
				}
			}
		}
	}
}

func MustExtractOwnerAndRepoFromGitUrl(repo string) (string, string) {
	u, err := url.Parse(repo)
	if err != nil {
		panic(err)
	}
	rawPath := strings.TrimPrefix(strings.TrimSuffix(u.Path, ".git"), "/")
	parts := strings.Split(rawPath, "/")
	return parts[0], parts[1]
}
