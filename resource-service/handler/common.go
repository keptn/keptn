package handler

import (
	"crypto/tls"
	"errors"
	git2go "github.com/libgit2/git2go/v33"
	"time"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common_models"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5/plumbing/transport/client"
	nethttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	ssh2 "golang.org/x/crypto/ssh"
)

const pathParamProjectName = "projectName"
const pathParamStageName = "stageName"
const pathParamServiceName = "serviceName"
const pathParamResourceURI = "resourceURI"

func OnAPIError(c *gin.Context, err error) {
	logger.Infof("Could not complete request %s %s: %v", c.Request.Method, c.Request.RequestURI, err)

	if check, resourceType := alreadyExists(err); check {
		SetConflictErrorResponse(c, resourceType+" already exists")
	} else if errors.Is(err, errors2.ErrProjectRepositoryNotEmpty) {
		SetConflictErrorResponse(c, "Project already exists with an already initialized GIT repository")
	} else if errors.Is(err, errors2.ErrInvalidGitToken) || errors.Is(err, errors2.ErrAuthenticationRequired) || errors.Is(err, errors2.ErrAuthorizationFailed) {
		SetFailedDependencyErrorResponse(c, "Invalid git token")
	} else if errors.Is(err, errors2.ErrCredentialsNotFound) {
		SetNotFoundErrorResponse(c, "Could not find credentials for upstream repository")
	} else if errors.Is(err, errors2.ErrMalformedCredentials) {
		SetFailedDependencyErrorResponse(c, "Could not decode credentials for upstream repository")
	} else if errors.Is(err, errors2.ErrCredentialsInvalidRemoteURL) || errors.Is(err, errors2.ErrCredentialsTokenMustNotBeEmpty) {
		SetBadRequestErrorResponse(c, "Upstream repository not found")
	} else if errors.Is(err, errors2.ErrRepositoryNotFound) {
		SetNotFoundErrorResponse(c, "Upstream repository not found")
	} else if check, resourceType := resourceNotFound(err); check {
		SetNotFoundErrorResponse(c, resourceType+" not found")
	} else {
		logger.Errorf("Encountered unknown error: %v", err)
		SetInternalServerErrorResponse(c, "Internal server error")
	}
}

func alreadyExists(err error) (bool, string) {
	if errors.Is(err, errors2.ErrProjectAlreadyExists) {
		return true, "Project"
	} else if errors.Is(err, errors2.ErrStageAlreadyExists) || errors.Is(err, errors2.ErrBranchExists) {
		return true, "Stage"
	} else if errors.Is(err, errors2.ErrServiceAlreadyExists) {
		return true, "Service"
	}
	return false, ""
}

func resourceNotFound(err error) (bool, string) {
	if errors.Is(err, errors2.ErrProjectNotFound) {
		return true, "Project"
	} else if errors.Is(err, errors2.ErrStageNotFound) || errors.Is(err, errors2.ErrReferenceNotFound) {
		return true, "Stage"
	} else if errors.Is(err, errors2.ErrServiceNotFound) {
		return true, "Service"
	} else if errors.Is(err, errors2.ErrResourceNotFound) {
		return true, "Resource"
	}
	return false, ""
}

func getAuthMethod(credentials *common_models.GitCredentials) (*common_models.AuthMethod, error) {
	if credentials.SshAuth != nil {
		return getSshGitAuth(credentials)
	} else if credentials.HttpsAuth != nil {
		return getHttpGitAuth(credentials)
	}
	return nil, nil
}

func getSshGitAuth(credentials *common_models.GitCredentials) (*common_models.AuthMethod, error) {
	publicKey, err := ssh.NewPublicKeys("git", []byte(credentials.SshAuth.PrivateKey), credentials.SshAuth.PrivateKeyPass)
	if err != nil {
		return nil, err
	}
	publicKey.HostKeyCallback = ssh2.InsecureIgnoreHostKey()

	credCallback := func(url string, usernameFromUrl string, allowedTypes git2go.CredentialType) (*git2go.Credential, error) {
		cred, err := git2go.NewCredentialSSHKeyFromMemory(credentials.User, string(publicKey.Signer.PublicKey().Marshal()), credentials.SshAuth.PrivateKey, credentials.SshAuth.PrivateKeyPass)
		if err != nil {
			return nil, err
		}
		return cred, nil
	}
	certCallback := func(cert *git2go.Certificate, valid bool, hostname string) error {
		return nil
	}

	return &common_models.AuthMethod{
		GoGitAuth: publicKey,
		Git2GoAuth: common_models.Git2GoAuth{
			CredCallback: credCallback,
			CertCallback: certCallback,
		},
	}, nil
}

func getHttpGitAuth(credentials *common_models.GitCredentials) (*common_models.AuthMethod, error) {
	git2GoAuth := common_models.Git2GoAuth{}
	if credentials.HttpsAuth.Proxy != nil {
		customClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: credentials.HttpsAuth.InsecureSkipTLS},
				Proxy: http.ProxyURL(&url.URL{
					Scheme: credentials.HttpsAuth.Proxy.Scheme,
					User:   url.UserPassword(credentials.HttpsAuth.Proxy.User, credentials.HttpsAuth.Proxy.Password),
					Host:   credentials.HttpsAuth.Proxy.URL,
				}),
			},

			// 15 second timeout
			Timeout: 15 * time.Second,

			// don't follow redirect
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		// Installing https protocol as a default one means that all the proxy traffic will be routed via secure connection
		// To use unsecure connection, InsecureSkipTLS parameter should be set to true and https protocol will be used without TLS verification
		client.InstallProtocol("https", nethttp.NewClient(customClient))

		git2GoAuth.ProxyOptions = git2go.ProxyOptions{
			Type: git2go.ProxyTypeAuto,
			Url:  credentials.HttpsAuth.Proxy.URL,
		}
	}

	if credentials.User == "" {
		//we try the authentication anyway since in most git servers
		//any user apart from an empty string is fine when we use a token
		//this auth will fail in case user is using bitbucket
		credentials.User = "keptnuser"
	}

	git2GoAuth.CredCallback = func(url string, usernameFromUrl string, allowedTypes git2go.CredentialType) (*git2go.Credential, error) {
		cred, err := git2go.NewCredentialUserpassPlaintext(credentials.User, credentials.HttpsAuth.Token)
		if err != nil {
			return nil, err
		}
		return cred, nil
	}

	return &common_models.AuthMethod{
		GoGitAuth: &nethttp.BasicAuth{
			Username: credentials.User,
			Password: credentials.HttpsAuth.Token,
		},
		Git2GoAuth: git2GoAuth,
	}, nil
}

func SetFailedDependencyErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusFailedDependency, models.Error{
		Code:    http.StatusFailedDependency,
		Message: msg,
	})
}

func SetNotFoundErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, models.Error{
		Code:    http.StatusNotFound,
		Message: msg,
	})
}

func SetInternalServerErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, models.Error{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
}

func SetBadRequestErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, models.Error{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
}

func SetConflictErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, models.Error{
		Code:    http.StatusConflict,
		Message: msg,
	})
}
