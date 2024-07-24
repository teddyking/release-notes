package kiln

type Kilnfile struct {
	Slug     string `yaml:"slug"`
	Releases []struct {
		Name             string `yaml:"name"`
		GithubRepository string `yaml:"github_repository"`
	} `yaml:"releases"`
}

func (k *Kilnfile) GetGithubRepositoryForRelease(name string) string {
	for _, r := range k.Releases {
		if r.Name == name {
			return r.GithubRepository
		}
	}
	return ""
}

type KilnfileLock struct {
	Releases []struct {
		Name         string `yaml:"name"`
		Sha1         string `yaml:"sha1"`
		Version      string `yaml:"version"`
		RemoteSource string `yaml:"remote_source"`
		RemotePath   string `yaml:"remote_path"`
	} `yaml:"releases"`
}

func (k *KilnfileLock) GetVersionForRelease(name string) string {
	for _, r := range k.Releases {
		if r.Name == name {
			return r.Version
		}
	}
	return ""
}
