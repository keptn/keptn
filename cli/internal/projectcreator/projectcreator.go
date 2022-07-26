package projectcreator

import (
	"encoding/base64"
	"errors"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"io"
	"io/fs"
	"strings"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/api_v1_interface.go . apiV1Interface:APIV1InterfaceMock
type apiV1Interface api.APIV1Interface

type EndpointInfo struct {
	APIEndpoint string
	APIToken    string
}

type ProjectInfo struct {
	Name              string
	Shipyard          string
	GitUser           string
	GitToken          string
	RemoteURL         string
	GitPrivateKey     string
	GitPrivateKeyPass string
	GitProxyURL       string
	GitProxyScheme    string
	GitProxyUser      string
	GitProxyPassword  string
	GitPemCertificate string
	InsecureSkipTLS   bool
}

type APIProvider func(EndpointInfo) (api.KeptnInterface, error)
type ShipyardProvider func(string) ([]byte, error)

type ProjectCreator struct {
	APIV1Interface   api.APIV1Interface
	ShipyardProvider ShipyardProvider
	FileSystem       fs.FS
}

func New(
	apiV1Interface api.APIV1Interface,
	shipyardProvider ShipyardProvider,
	fileSystem fs.FS) *ProjectCreator {
	return &ProjectCreator{
		APIV1Interface:   apiV1Interface,
		ShipyardProvider: shipyardProvider,
		FileSystem:       fileSystem,
	}
}

// CreateProject does what is required to create a new Keptn Project
// First it gets the metadata from the Keptn API to see whether the Keptn installation was setup
// using the automatic provisioning feature, then it tries to load the shipyard file and finally
// it tries to create the Keptn project using the public API of Keptn
func (p *ProjectCreator) CreateProject(projectInfo ProjectInfo) error {
	metaData, errMetadata := p.APIV1Interface.GetMetadata()
	if errMetadata != nil {
		return fmt.Errorf("Unable to fetch metadata: %s\n", errMetadata.GetMessage())
	}

	shipyard, errGettingEndpoint := p.ShipyardProvider(projectInfo.Shipyard)
	if errGettingEndpoint != nil {
		return fmt.Errorf("Failed to read and parse shipyard file - %s", errGettingEndpoint.Error())
	}

	encodedShipyardContent := base64.StdEncoding.EncodeToString(shipyard)
	project := apimodels.CreateProject{
		Name:     &projectInfo.Name,
		Shipyard: &encodedShipyardContent,
	}

	if !*metaData.Automaticprovisioning {
		if strNotSet(projectInfo.GitUser) || strNotSet(projectInfo.RemoteURL) {
			return errors.New("GIT username and GIT remote URL must be set")
		}
		if strNotSet(projectInfo.GitToken) && strNotSet(projectInfo.GitPrivateKey) {
			return errors.New("GIT Access token or GIT private key must be set")
		}

		if strSet(projectInfo.GitToken) && strSet(projectInfo.GitPrivateKey) {
			return errors.New("GIT Access token and GIT private key cannot be set together")
		}

		if strSet(projectInfo.GitProxyURL) && strings.HasPrefix(projectInfo.RemoteURL, "ssh://") {
			return errors.New("GIT Proxy cannot be set with SSH")
		}

		if strSet(projectInfo.GitProxyURL) && strNotSet(projectInfo.GitProxyScheme) {
			return errors.New("GIT Proxy cannot be set without scheme")
		}

		project.GitCredentials = &apimodels.GitAuthCredentials{
			User:      projectInfo.GitUser,
			RemoteURL: projectInfo.RemoteURL,
		}

		if strings.HasPrefix(projectInfo.RemoteURL, "ssh://") {
			privateKeyFile, err := p.FileSystem.Open(projectInfo.GitPrivateKey)
			if err != nil {
				return err
			}
			privateKeyFileContent, err := io.ReadAll(privateKeyFile)
			defer privateKeyFile.Close()
			if err != nil {
				return fmt.Errorf("unable to read privateKey file: %s\n", err.Error())
			}

			sshCredentials := apimodels.SshGitAuth{
				PrivateKey:     base64.StdEncoding.EncodeToString(privateKeyFileContent),
				PrivateKeyPass: projectInfo.GitPrivateKeyPass,
			}

			project.GitCredentials.SshAuth = &sshCredentials
		} else if strings.HasPrefix(projectInfo.RemoteURL, "http") {
			httpCredentials := apimodels.HttpsGitAuth{
				Token:           projectInfo.GitToken,
				InsecureSkipTLS: projectInfo.InsecureSkipTLS,
			}
			if strSet(projectInfo.GitProxyURL) {
				proxyCredentials := apimodels.ProxyGitAuth{
					URL:      projectInfo.GitProxyURL,
					Scheme:   projectInfo.GitProxyScheme,
					User:     projectInfo.GitProxyUser,
					Password: projectInfo.GitProxyPassword,
				}
				httpCredentials.Proxy = &proxyCredentials
			}

			if strSet(projectInfo.GitPemCertificate) {
				gitCertFile, err := p.FileSystem.Open(projectInfo.GitPemCertificate)

				if err != nil {
					return err
				}
				gitCertFileContent, err := io.ReadAll(gitCertFile)
				defer gitCertFile.Close()
				if err != nil {
					return err
				}
				httpCredentials.Certificate = base64.StdEncoding.EncodeToString(gitCertFileContent)
			}
			project.GitCredentials.HttpsAuth = &httpCredentials
		}
	}
	_, errCreateProject := p.APIV1Interface.CreateProject(project)
	if errCreateProject != nil {
		return fmt.Errorf("Create project was unsuccessful.\n%s", errCreateProject.GetMessage())
	}
	logging.PrintLog("Project created successfully", logging.InfoLevel)
	return nil

}

func strSet(s string) bool {
	return s != ""
}

func strNotSet(s string) bool {
	return s == ""
}
