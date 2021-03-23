package cmd

import (
	"github.com/bmizerany/assert"
	"github.com/keptn/go-utils/pkg/api/models"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	fakeapi "github.com/keptn/go-utils/pkg/api/utils/fake"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	credentialmanager_mock "github.com/keptn/keptn/cli/pkg/credentialmanager/fake"
	"net/url"
	"testing"
)

func TestSecretCmdHandler_CreateSecret(t *testing.T) {
	type fields struct {
		credentialManager credentialmanager.CredentialManagerInterface
		secretAPI         *fakeapi.SecretHandlerInterfaceMock
	}
	type args struct {
		secretName string
		data       []string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantSecretObj *apimodels.Secret
		wantErr       bool
	}{
		{
			name: "create secret",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					CreateSecretFunc: func(secret apimodels.Secret) (string, *apimodels.Error) {
						return "", nil
					},
				},
			},
			args: args{
				secretName: "my-secret",
				data:       []string{"foo=bar", "bar=foo"},
			},
			wantSecretObj: &apimodels.Secret{
				Data: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
				Name:  stringp("my-secret"),
				Scope: stringp(defaultSecretScope),
			},
			wantErr: false,
		},
		{
			name: "create secret failed - return error",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					CreateSecretFunc: func(secret apimodels.Secret) (string, *apimodels.Error) {
						return "", &apimodels.Error{Message: stringp("could not update secret")}
					},
				},
			},
			args: args{
				secretName: "my-secret",
				data:       []string{"foo=bar", "bar=foo"},
			},
			wantSecretObj: &apimodels.Secret{
				Data: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
				Name:  stringp("my-secret"),
				Scope: stringp(defaultSecretScope),
			},
			wantErr: true,
		},
		{
			name: "invalid secret - return error",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					CreateSecretFunc: func(secret apimodels.Secret) (string, *apimodels.Error) {
						return "", &apimodels.Error{Message: stringp("could not update secret")}
					},
				},
			},
			args: args{
				secretName: "my-secret",
				data:       []string{"foobar", "bar=foo"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := SecretCmdHandler{
				credentialManager: tt.fields.credentialManager,
				secretAPI:         tt.fields.secretAPI,
			}
			if err := h.CreateSecret(tt.args.secretName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("CreateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantSecretObj != nil {
				assert.Equal(t, 1, len(tt.fields.secretAPI.CreateSecretCalls()))

				updateCall := tt.fields.secretAPI.CreateSecretCalls()[0]
				assert.Equal(t, *tt.wantSecretObj, updateCall.Secret)
			} else {
				assert.Equal(t, 0, len(tt.fields.secretAPI.CreateSecretCalls()))
			}
		})
	}
}

func TestSecretCmdHandler_UpdateSecret(t *testing.T) {
	type fields struct {
		credentialManager credentialmanager.CredentialManagerInterface
		secretAPI         *fakeapi.SecretHandlerInterfaceMock
	}
	type args struct {
		secretName string
		data       []string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantSecretObj *apimodels.Secret
		wantErr       bool
	}{
		{
			name: "update secret",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					UpdateSecretFunc: func(secret apimodels.Secret) (string, *apimodels.Error) {
						return "", nil
					},
				},
			},
			args: args{
				secretName: "my-secret",
				data:       []string{"foo=bar", "bar=foo"},
			},
			wantSecretObj: &apimodels.Secret{
				Data: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
				Name:  stringp("my-secret"),
				Scope: stringp(defaultSecretScope),
			},
			wantErr: false,
		},
		{
			name: "update secret failed -return error",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					UpdateSecretFunc: func(secret apimodels.Secret) (string, *apimodels.Error) {
						return "", &apimodels.Error{Message: stringp("could not update secret")}
					},
				},
			},
			args: args{
				secretName: "my-secret",
				data:       []string{"foo=bar", "bar=foo"},
			},
			wantSecretObj: &apimodels.Secret{
				Data: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
				Name:  stringp("my-secret"),
				Scope: stringp(defaultSecretScope),
			},
			wantErr: true,
		},
		{
			name: "invalid secret - return error",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					UpdateSecretFunc: func(secret apimodels.Secret) (string, *apimodels.Error) {
						return "", &apimodels.Error{Message: stringp("could not update secret")}
					},
				},
			},
			args: args{
				secretName: "my-secret",
				data:       []string{"foobar", "bar=foo"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := SecretCmdHandler{
				credentialManager: tt.fields.credentialManager,
				secretAPI:         tt.fields.secretAPI,
			}
			if err := h.UpdateSecret(tt.args.secretName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantSecretObj != nil {
				assert.Equal(t, 1, len(tt.fields.secretAPI.UpdateSecretCalls()))

				updateCall := tt.fields.secretAPI.UpdateSecretCalls()[0]
				assert.Equal(t, *tt.wantSecretObj, updateCall.Secret)
			}

		})
	}
}

func TestSecretCmdHandler_DeleteSecret(t *testing.T) {
	type fields struct {
		credentialManager credentialmanager.CredentialManagerInterface
		secretAPI         api.SecretHandlerInterface
	}
	type args struct {
		name  string
		scope string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete secret",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					DeleteSecretFunc: func(secretName string, secretScope string) (string, *models.Error) {
						return "", nil
					},
				},
			},
			args: args{
				name:  "my-secret",
				scope: defaultSecretScope,
			},
			wantErr: false,
		},
		{
			name: "delete secret failed - return error",
			fields: fields{
				credentialManager: createMockCredentialManager(),
				secretAPI: &fakeapi.SecretHandlerInterfaceMock{
					DeleteSecretFunc: func(secretName string, secretScope string) (string, *models.Error) {
						return "", &apimodels.Error{
							Code:    0,
							Message: stringp("could not delete secret"),
						}
					},
				},
			},
			args: args{
				name:  "my-secret",
				scope: defaultSecretScope,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := SecretCmdHandler{
				credentialManager: tt.fields.credentialManager,
				secretAPI:         tt.fields.secretAPI,
			}
			if err := h.DeleteSecret(tt.args.name, tt.args.scope); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createMockCredentialManager() *credentialmanager_mock.CredentialManagerInterfaceMock {
	return &credentialmanager_mock.CredentialManagerInterfaceMock{
		GetCredsFunc: func(namespace string) (url.URL, string, error) {
			return url.URL{}, "", nil
		},
	}
}
