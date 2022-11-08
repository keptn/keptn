package common

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	git2go "github.com/libgit2/git2go/v33"
	log "github.com/sirupsen/logrus"
	"strings"
)

const unexpectedClientErrorMessage = "unexpected client"

//go:generate moq -pkg common_mock -skip-ensure -out ./fake/gogit_mock.go . Gogit
type Gogit interface {
	PlainOpen(path string) (*git.Repository, error)
	PlainClone(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	PlainInit(path string, isBare bool) (*git.Repository, error)
}

type GogitReal struct{}

func (t GogitReal) PlainOpen(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func (t GogitReal) PlainClone(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	repo, err := git.PlainClone(path, isBare, o)
	if err != nil {
		log.Warnf("Could not clone using go-git library: %v", err)
		if strings.Contains(err.Error(), unexpectedClientErrorMessage) {
			git2 := Git2Go{}
			log.Debug("Try to clone using libgit2go")
			return git2.PlainClone(gitContext, path, isBare, o)
		}
		return nil, err
	}
	return repo, nil
}

func (t GogitReal) PlainInit(path string, isBare bool) (*git.Repository, error) {
	return git.PlainInit(path, isBare)
}

type Git2Go struct{}

func (g Git2Go) PlainClone(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	_, err := git2go.Clone(o.URL, path, &git2go.CloneOptions{
		Bare: isBare,
		FetchOptions: git2go.FetchOptions{
			DownloadTags: git2go.DownloadTagsNone,
			RemoteCallbacks: git2go.RemoteCallbacks{
				CredentialsCallback:      gitContext.AuthMethod.Git2GoAuth.CredCallback,
				CertificateCheckCallback: gitContext.AuthMethod.Git2GoAuth.CertCallback,
			},
			ProxyOptions: gitContext.AuthMethod.Git2GoAuth.ProxyOptions,
		},
	})

	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	_, err = repo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return nil, kerrors.ErrEmptyRemoteRepository
		}
		return nil, err
	}
	log.Debug("Clone using libgit2go succeeded")
	return repo, nil
}
