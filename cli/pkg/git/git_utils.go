package git

import (
	"github.com/go-git/go-git/v5"
)

func CloneGitHubUrl(filePath string, url string) error {
	_, err := git.PlainClone(filePath, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return err
	}
	return nil
}
