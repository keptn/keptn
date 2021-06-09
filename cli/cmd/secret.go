package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"gopkg.in/yaml.v3"
	"strings"
)

func parseSecretData(in []string) (map[string]string, error) {
	result := map[string]string{}

	for _, literal := range in {
		key, value, err := parseLiteralKeyValuePair(literal)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

const defaultSecretScope = "keptn-default"

type SecretCmdHandler struct {
	credentialManager credentialmanager.CredentialManagerInterface
	secretAPI         api.SecretHandlerInterface
}

func (h SecretCmdHandler) CreateSecret(secretName string, data []string, scope *string) error {
	var secretScope string
	if scope == nil || *scope == "" {
		secretScope = defaultSecretScope
	} else {
		secretScope = *scope
	}
	secretData, err := parseSecretData(data)
	if err != nil {
		return err
	}
	if _, err := h.secretAPI.CreateSecret(models.Secret{
		Data: secretData,
		SecretMetadata: models.SecretMetadata{
			Name:  &secretName,
			Scope: &secretScope,
		},
	}); err != nil {
		return errors.New(*err.Message)
	}
	return nil
}

func (h SecretCmdHandler) UpdateSecret(secretName string, data []string, scope *string) error {
	var secretScope string
	if scope == nil || *scope == "" {
		secretScope = defaultSecretScope
	} else {
		secretScope = *scope
	}
	secretData, err := parseSecretData(data)
	if err != nil {
		return err
	}
	secret := models.Secret{
		Data: secretData,
		SecretMetadata: models.SecretMetadata{
			Name:  &secretName,
			Scope: &secretScope,
		},
	}
	if _, err := h.secretAPI.UpdateSecret(secret); err != nil {
		return errors.New(*err.Message)
	}
	return nil
}

func (h SecretCmdHandler) DeleteSecret(name string, scope *string) error {
	var secretScope string
	if scope == nil || *scope == "" {
		secretScope = defaultSecretScope
	} else {
		secretScope = *scope
	}
	if _, err := h.secretAPI.DeleteSecret(name, secretScope); err != nil {
		return errors.New(*err.Message)
	}
	return nil
}

func (h SecretCmdHandler) GetSecrets(outputFormat string) (string, error) {
	secrets, errObj := h.secretAPI.GetSecrets()
	if errObj != nil {
		return "", errors.New(*errObj.Message)
	}

	var output string
	if outputFormat == "json" {
		marshal, err := json.MarshalIndent(secrets, "", "  ")
		if err != nil {
			return "", err
		}
		output = string(marshal)
	} else if outputFormat == "yaml" {
		marshal, err := yaml.Marshal(secrets)
		if err != nil {
			return "", err
		}
		output = string(marshal)
	} else {
		if len(secrets.Secrets) == 0 {
			output = "No secrets found"
		} else {
			output = "NAME"
			for _, secret := range secrets.Secrets {
				output = output + "\n" + *secret.Name
			}
		}
	}
	return output, nil
}

func parseLiteralKeyValuePair(in string) (string, string, error) {
	// leading equal is invalid
	if strings.Index(in, "=") == 0 {
		return "", "", fmt.Errorf("invalid literal source %v, expected key=value", in)
	}
	// split after the first equal
	split := strings.SplitN(in, "=", 2)
	if len(split) != 2 {
		return "", "", fmt.Errorf("invalid literal source %v, expected key=value", in)
	}

	return split[0], split[1], nil
}

func NewSecretCmdHandler(cm credentialmanager.CredentialManagerInterface) (*SecretCmdHandler, error) {
	sh := &SecretCmdHandler{credentialManager: cm}
	endPoint, apiToken, err := cm.GetCreds(namespace)
	if err != nil {
		return nil, errors.New(authErrorMsg)
	}
	if endPointErr := CheckEndpointStatus(endPoint.String()); endPointErr != nil {
		return nil, fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
			endPointErr)
	}
	sh.secretAPI = api.NewAuthenticatedSecretHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
	return sh, nil
}
