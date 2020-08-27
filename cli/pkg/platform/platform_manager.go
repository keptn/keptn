package platform

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/file"
)

const OpenShiftIdentifier = "openshift"
const KubernetesIdentifier = "kubernetes"

type platform interface {
	checkRequirements() error
	getCreds() interface{}
	checkCreds() error
	readCreds()
	printCreds()
}

type PlatformManager struct {
	platform platform
}

func NewPlatformManager(platformIdentifier string) (*PlatformManager, error) {

	switch strings.ToLower(platformIdentifier) {
	case OpenShiftIdentifier:
		return &PlatformManager{platform: newOpenShiftPlatform()}, nil
	case KubernetesIdentifier:
		return &PlatformManager{platform: newKubernetesPlatform()}, nil
	default:
		return nil, errors.New("Unsupported platform '" + platformIdentifier +
			"'. The following platforms are supported: OpenShiftIdentifier and KubernetesIdentifier")
	}
}

func (mng PlatformManager) CheckRequirements() error {
	return mng.platform.checkRequirements()
}

func (mng PlatformManager) CheckCreds() error {
	return mng.platform.checkCreds()
}

func (mng PlatformManager) ParseConfig(configFile string) error {
	data, err := file.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), mng.platform.getCreds())
}

func (mng PlatformManager) ReadCreds() error {

	cm := credentialmanager.NewCredentialManager()
	credsStr, err := cm.GetInstallCreds()
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
	}

	newCreds, _ := json.Marshal(mng.platform.getCreds())
	newCredsStr := strings.Replace(string(newCreds), "\r\n", "\n", -1)
	newCredsStr = strings.Replace(newCredsStr, "\n", "", -1)
	return cm.SetInstallCreds(newCredsStr)
}
