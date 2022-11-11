package common

import (
	"github.com/go-git/go-git/v5"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestGit2Go_PlainClone(t *testing.T) {
	defer func() {
		_ = os.RemoveAll(TESTPATH + "/upstreamrepo")
		_ = os.RemoveAll(TESTPATH + "/localrepo")
	}()
	// make local remote
	upstreamUrl := TESTPATH + "/upstreamrepo"
	_, err := git.PlainClone(upstreamUrl, true, &git.CloneOptions{URL: "https://github.com/git-fixtures/basic.git"})
	require.Nil(t, err)

	g := Git2Go{}

	httpCredentials := apimodels.HttpsGitAuth{
		Token: "mytoken",
	}
	gitContext := common_models.GitContext{Project: "localrepo",
		Credentials: &common_models.GitCredentials{
			User:      "u2",
			HttpsAuth: &httpCredentials,
			RemoteURL: upstreamUrl,
		}}

	repo, err := g.PlainClone(gitContext, TESTPATH+"/localrepo", true, &git.CloneOptions{
		URL:             gitContext.Credentials.RemoteURL,
		Auth:            gitContext.AuthMethod.GoGitAuth,
		InsecureSkipTLS: retrieveInsecureSkipTLS(gitContext.Credentials),
	})

	require.Nil(t, err)
	require.NotNil(t, repo)
}

func TestGit2Go_PlainClone_EmptyUpstream(t *testing.T) {
	defer func() {
		_ = os.RemoveAll(TESTPATH + "/upstreamrepo")
		_ = os.RemoveAll(TESTPATH + "/localrepo")
	}()
	// make empty local remote
	upstreamUrl := TESTPATH + "/upstreamrepo"
	emptyUpstreamUrl, err := filepath.Abs(upstreamUrl)
	require.Nil(t, err)

	_, err = git.PlainInit(emptyUpstreamUrl, true)
	require.Nil(t, err)

	g := Git2Go{}

	httpCredentials := apimodels.HttpsGitAuth{
		Token: "mytoken",
	}
	gitContext := common_models.GitContext{Project: "localrepo",
		Credentials: &common_models.GitCredentials{
			User:      "u2",
			HttpsAuth: &httpCredentials,
			RemoteURL: emptyUpstreamUrl,
		}}

	repo, err := g.PlainClone(gitContext, TESTPATH+"/localrepo", true, &git.CloneOptions{
		URL:             gitContext.Credentials.RemoteURL,
		Auth:            gitContext.AuthMethod.GoGitAuth,
		InsecureSkipTLS: retrieveInsecureSkipTLS(gitContext.Credentials),
	})

	require.ErrorIs(t, err, errors.ErrEmptyRemoteRepository)
	require.Nil(t, repo)
}
