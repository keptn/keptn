package cmd

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
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
		Data:  secretData,
		Name:  &secretName,
		Scope: &secretScope,
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
		Data:  secretData,
		Name:  &secretName,
		Scope: &secretScope,
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
