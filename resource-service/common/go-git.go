package common

import (
	"github.com/go-git/go-git/v5"
)

//go:generate moq -pkg common_mock -skip-ensure -out ./fake/gogit_mock.go . Gogit
type Gogit interface {
	PlainOpen(path string) (*git.Repository, error)
	PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	PlainInit(path string, isBare bool) (*git.Repository, error)
}

type GogitReal struct{}

func (t GogitReal) PlainOpen(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func (t GogitReal) PlainClone(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, isBare, o)
}

func (t GogitReal) PlainInit(path string, isBare bool) (*git.Repository, error) {
	return git.PlainInit(path, isBare)
}

/*//go:generate moq -pkg common_mock -skip-ensure -out ./fake/gogitrepo_mock.go . Repository
type Repository interface {
	Remote(name string) (*git.Remote, error)
	Remotes() ([]*git.Remote, error)
	CreateRemote(c *config.RemoteConfig) (*git.Remote, error)
	CreateRemoteAnonymous(c *config.RemoteConfig) (*git.Remote, error)
	DeleteRemote(name string) error
	Branch(name string) (*config.Branch, error)
	CreateBranch(c *config.Branch) error
	DeleteBranch(name string) error
	resolveToCommitHash(h plumbing.Hash) (plumbing.Hash, error)
	clone(ctx context.Context, o *git.CloneOptions) error
	calculateRemoteHeadReference(spec []config.RefSpec, resolvedHead *plumbing.Reference) []*plumbing.Reference
	Push(o *git.PushOptions) error
	Branches() (storer.ReferenceIter, error)
	TreeObject(h plumbing.Hash) (*object.Tree, error)
	TreeObjects() (*object.TreeIter, error)
	CommitObject(h plumbing.Hash) (*object.Commit, error)
	CommitObjects() (object.CommitIter, error)
	BlobObject(h plumbing.Hash) (*object.Blob, error)
	BlobObjects() (*object.BlobIter, error)
	Object(t plumbing.ObjectType, h plumbing.Hash) (object.Object, error)
	Objects() (*object.ObjectIter, error)
	Head() (*plumbing.Reference, error)
	Reference(name plumbing.ReferenceName, resolved bool) (*plumbing.Reference, error)
	References() (storer.ReferenceIter, error)
	Worktree() (*git.Worktree, error)
	ResolveRevision(rev plumbing.Revision) (*plumbing.Hash, error)
}*/
