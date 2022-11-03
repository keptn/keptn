package common

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/keptn/keptn/resource-service/common_models"
	git2go "github.com/libgit2/git2go/v34"
)

//go:generate moq -pkg common_mock -skip-ensure -out ./fake/gogit_mock.go . Gogit
type Gogit interface {
	PlainOpen(path string) (*git.Repository, error)
	PlainClone(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	PlainInit(gitContext common_models.GitContext, path string, isBare bool) (*git.Repository, error)
	Push(gitContext common_models.GitContext, repository *git.Repository, options *git.PushOptions) error
	Pull(gitContext common_models.GitContext, worktree *git.Worktree, options *git.PullOptions) error
	Fetch(gitContext common_models.GitContext, repository *git.Repository, options *git.FetchOptions) error
}

type GogitReal struct{}

func (t GogitReal) PlainOpen(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func (t GogitReal) PlainClone(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, isBare, o)
}

func (t GogitReal) PlainInit(gitContext common_models.GitContext, path string, isBare bool) (*git.Repository, error) {
	return git.PlainInit(path, isBare)
}

func (t GogitReal) Push(gitContext common_models.GitContext, repository *git.Repository, o *git.PushOptions) error {
	return repository.Push(o)
}

func (t GogitReal) Pull(gitContext common_models.GitContext, workTree *git.Worktree, o *git.PullOptions) error {
	return workTree.Pull(o)
}

func (t GogitReal) Fetch(gitContext common_models.GitContext, repository *git.Repository, o *git.FetchOptions) error {
	return repository.Fetch(o)
}

type Git2Go struct{}

func (g Git2Go) PlainOpen(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func (g Git2Go) PlainClone(gitContext common_models.GitContext, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	_, err := git2go.Clone(o.URL, path, &git2go.CloneOptions{
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

	return git.PlainOpen(path)
}

func (g Git2Go) PlainInit(gitContext common_models.GitContext, path string, isBare bool) (*git.Repository, error) {
	return git.PlainInit(path, isBare)
}

func (g Git2Go) Push(gitContext common_models.GitContext, repository *git.Repository, options *git.PushOptions) error {
	projectConfigPath := GetProjectConfigPath(gitContext.Project)

	repo, err := git2go.OpenRepository(projectConfigPath)
	if err != nil {
		return err
	}

	head, err := repository.Head()
	if err != nil {
		return err
	}

	remote, err := repo.Remotes.Lookup(options.RemoteName)
	if err != nil {
		return nil
	}

	err = remote.Push([]string{head.Name().String()}, &git2go.PushOptions{
		RemoteCallbacks: git2go.RemoteCallbacks{
			CredentialsCallback:      gitContext.AuthMethod.Git2GoAuth.CredCallback,
			CertificateCheckCallback: gitContext.AuthMethod.Git2GoAuth.CertCallback,
		},
		ProxyOptions: gitContext.AuthMethod.Git2GoAuth.ProxyOptions,
	})

	// TODO error mapping
	return err
}

func (g Git2Go) Pull(gitContext common_models.GitContext, worktree *git.Worktree, options *git.PullOptions) error {
	projectConfigPath := GetProjectConfigPath(gitContext.Project)

	repo, err := git2go.OpenRepository(projectConfigPath)
	if err != nil {
		return err
	}

	// Locate remote
	remote, err := repo.Remotes.Lookup(options.RemoteName)
	if err != nil {
		return err
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, &git2go.FetchOptions{
		RemoteCallbacks: git2go.RemoteCallbacks{
			CredentialsCallback:      gitContext.AuthMethod.Git2GoAuth.CredCallback,
			CertificateCheckCallback: gitContext.AuthMethod.Git2GoAuth.CertCallback,
		},
		ProxyOptions: gitContext.AuthMethod.Git2GoAuth.ProxyOptions,
	}, ""); err != nil {
		return err
	}

	// Get remote ref
	head, err := repo.Head()
	if err != nil {
		return err
	}
	refName := fmt.Sprintf("refs/remotes/%s/%s", options.RemoteName, head.Branch().Reference.Shorthand())
	remoteBranch, err := repo.References.Lookup(refName)
	if err != nil {
		return err
	}

	remoteBranchID := remoteBranch.Target()
	// Get annotated commit
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		return err
	}

	// Do the merge analysis
	mergeHeads := make([]*git2go.AnnotatedCommit, 1)
	mergeHeads[0] = annotatedCommit
	analysis, _, err := repo.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}

	if analysis&git2go.MergeAnalysisUpToDate != 0 {
		return nil
	} else if analysis&git2go.MergeAnalysisNormal != 0 {
		// Just merge changes
		if err := repo.Merge([]*git2go.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
			return err
		}
		// Check for conflicts
		index, err := repo.Index()
		if err != nil {
			return err
		}

		if index.HasConflicts() {
			return git.ErrNonFastForwardUpdate
		}

		// Make the merge commit
		sig, err := repo.DefaultSignature()
		if err != nil {
			return err
		}

		// Get Write Tree
		treeId, err := index.WriteTree()
		if err != nil {
			return err
		}

		tree, err := repo.LookupTree(treeId)
		if err != nil {
			return err
		}

		localCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}

		remoteCommit, err := repo.LookupCommit(remoteBranchID)
		if err != nil {
			return err
		}

		repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)

		// Clean up
		repo.StateCleanup()
	} else if analysis&git2go.MergeAnalysisFastForward != 0 {
		// Fast-forward changes
		// Get remote tree
		remoteTree, err := repo.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}

		// Checkout
		if err := repo.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}

		branchRef, err := repo.References.Lookup("refs/heads/" + head.Branch().Reference.Shorthand())
		if err != nil {
			return err
		}

		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("unexpected merge analysis result %d", analysis)
	}

	return nil
}

func (g Git2Go) Fetch(gitContext common_models.GitContext, repository *git.Repository, options *git.FetchOptions) error {
	projectConfigPath := GetProjectConfigPath(gitContext.Project)

	repo, err := git2go.OpenRepository(projectConfigPath)
	if err != nil {
		return err
	}

	// Locate remote
	remote, err := repo.Remotes.Lookup(options.RemoteName)
	if err != nil {
		return err
	}

	// Fetch changes from remote
	return remote.Fetch([]string{}, &git2go.FetchOptions{
		RemoteCallbacks: git2go.RemoteCallbacks{
			CredentialsCallback:      gitContext.AuthMethod.Git2GoAuth.CredCallback,
			CertificateCheckCallback: gitContext.AuthMethod.Git2GoAuth.CertCallback,
		},
		ProxyOptions: gitContext.AuthMethod.Git2GoAuth.ProxyOptions,
	}, "")
}
