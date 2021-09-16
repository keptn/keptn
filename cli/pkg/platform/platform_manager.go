package platform

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/fileutils"
	"os"
	"strings"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
)

// OpenShiftIdentifier is used as identifier for OpenShift
const OpenShiftIdentifier = "openshift"

// KubernetesIdentifier is used as identifier for Kubernetes
const KubernetesIdentifier = "kubernetes"

type platform interface {
	checkRequirements() error
	getCreds() interface{}
	checkCreds() error
	readCreds()
	printCreds()
}

// PlatformManager allows to manage (i.e. check requirements) for supported platforms
type PlatformManager struct {
	platform          platform
	credentialManager credentialmanager.CredentialManagerInterface
}

// NewPlatformManager creates a new manager for the provided platform
func NewPlatformManager(platformIdentifier string, cm credentialmanager.CredentialManagerInterface) (*PlatformManager, error) {

	switch strings.ToLower(platformIdentifier) {
	case OpenShiftIdentifier:
		return &PlatformManager{platform: newOpenShiftPlatform(), credentialManager: cm}, nil
	case KubernetesIdentifier:
		return &PlatformManager{platform: newKubernetesPlatform(), credentialManager: cm}, nil
	default:
		return nil, errors.New("Unsupported platform '" + platformIdentifier +
			"'. The following platforms are supported: OpenShiftIdentifier and KubernetesIdentifier")
	}
}

func (mng *PlatformManager) GetCredentialManager() credentialmanager.CredentialManagerInterface {
	return mng.credentialManager
}

// CheckRequirements checks the platform's requirements
func (mng *PlatformManager) CheckRequirements() error {
	return mng.platform.checkRequirements()
}

// CheckCreds checks the provided creds
func (mng *PlatformManager) CheckCreds() error {
	return mng.platform.checkCreds()
}

// ParseConfig reads and parses the provided config file
func (mng *PlatformManager) ParseConfig(configFile string) error {
	data, err := fileutils.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, mng.platform.getCreds())
}

// ReadCreds reads the credentials for the platform
func (mng *PlatformManager) ReadCreds(assumeYes bool) error {
	credsStr, err := mng.credentialManager.GetInstallCreds()
	if err != nil {
		credsStr = ""
	}
	// Ignore unmarshaling error
	json.Unmarshal([]byte(credsStr), mng.platform.getCreds())

	for {
		mng.platform.readCreds()

		fmt.Println()
		fmt.Println("Please confirm that the provided cluster information is correct: ")

		mng.platform.printCreds()
		if assumeYes {
			break
		}
		fmt.Println("Is this all correct? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.TrimSpace(in)
		if in == "y" || in == "yes" {
			break
		}

		if in == "n" || in == "no" {
			return errors.New("Stopping installation")
		}
	}

	newCreds, _ := json.Marshal(mng.platform.getCreds())
	newCredsStr := strings.Replace(string(newCreds), "\r\n", "\n", -1)
	newCredsStr = strings.Replace(newCredsStr, "\n", "", -1)
	return mng.credentialManager.SetInstallCreds(newCredsStr)
}
